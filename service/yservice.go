package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/api/youtube/v3"

	"github.com/gen1us2k/log"
	"github.com/maddevsio/yourcast-streamer/bot"
	"github.com/maddevsio/yourcast-streamer/service/data"
	"github.com/maddevsio/yourcast-streamer/stream"
)

// YoutubeStreamService re-streams youtube video
// on the fly to rtmp server
type YoutubeStreamService struct {
	BaseService

	s  *Streamer
	ss *data.StreamStorage
	yc *bot.YoutubeClient

	logger log.Logger
}

// Name returns name of service
func (ys *YoutubeStreamService) Name() string {
	return "youtube_stream_service"
}

// Init initializes logger and service
func (ys *YoutubeStreamService) Init(s *Streamer) error {
	ys.s = s
	ys.logger = log.NewLogger(ys.Name())
	yc, err := bot.NewYoutubeClient(ys.s.Config().YoutubeAPIKey)
	if err != nil {
		return err
	}
	ys.ss = data.NewStreamStorage()
	ys.yc = yc
	return nil
}

// Run runs YoutubeStreamService
func (ys *YoutubeStreamService) Run() error {
	ys.logger.Info("Getting current streams")
	streams, err := ys.getStreams()
	if err != nil {
		return err
	}
	ys.logger.Info("Streams received. Populating internal storage")
	for _, stream := range streams {
		if stream.IsAutoStream() {
			ys.runJobsForAutoStream(stream, false)
			continue
		}

		ys.ss.Lock()
		ys.ss.Items[stream.ID] = stream.ToStreamItem()
		ys.ss.Unlock()
	}
	for _, item := range ys.ss.Items {
		ys.logger.Infof("Starting streaming of %s", item.Name)
		if !ys.s.Config().DisableStreaming {
			ys.s.waitGroup.Add(1)
			go ys.runStream(item)
		}
		if !item.IsAuto {
			ys.s.waitGroup.Add(1)
			go ys.downloadStream(item)
		}
	}
	return nil
}

func (ys *YoutubeStreamService) getStreams() ([]data.Stream, error) {
	url := fmt.Sprintf("%s/api/streams/", ys.s.Config().WebUIURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var streams []data.Stream
	err = json.Unmarshal(body, &streams)
	if err != nil {
		return nil, err
	}
	return streams, nil
}

func (ys *YoutubeStreamService) runStream(data data.StreamItem) {
	data.Lock()
	e := data.Links.Front()
	data.Unlock()
	ys.logger.Infof("Preparing to stream items in %s channel", data.Name)
	for {
		if e == nil {
			continue
		}
		youtubeURL := fmt.Sprintf("%v", e.Value)
		absFileName := stream.GetFileNameByURL(youtubeURL, ys.s.Config().RootPath)
		dstURL := fmt.Sprintf("%s/%s", ys.s.Config().RTMPRootServerURL, data.Slug)
		if _, err := os.Stat(absFileName); err == nil {
			ys.logger.Infof(
				"Streaming channel %s video %s from file",
				data.Name, youtubeURL,
			)
			err := stream.FromLocalFile(
				ys.s.Config().FFMpegPath,
				absFileName, dstURL,
			)
			if err != nil {
				ys.logger.Errorf("Got error %s while streaming video %s for channel %s", err, youtubeURL, data.Name)
			}
		} else {
			ys.logger.Infof(
				"Streaming channel %s video %s from Youtube",
				data.Name, youtubeURL,
			)
			err := stream.FromYoutubeURL(
				ys.s.Config().FFMpegPath,
				youtubeURL, dstURL, ys.s.Config().RootPath,
			)
			if err != nil {
				ys.logger.Errorf("Got error %s while streaming video %s for channel %s", err, youtubeURL, data.Name)
			}
		}
		data.Lock()
		e = e.Next()
		if e == nil {
			e = data.Links.Front()
		}
		data.Unlock()
	}
	ys.logger.Infof("Streaming of %s channel done", data.Name)
	ys.s.waitGroup.Done()
}

func (ys *YoutubeStreamService) downloadStream(data data.StreamItem) {
	for e := data.Links.Front(); e != nil; e = e.Next() {

		youtubeURL := fmt.Sprintf("%v", e.Value)
		absFileName := stream.GetFileNameByURL(youtubeURL, ys.s.Config().RootPath)
		if !stream.FileExist(absFileName) {

			ys.logger.Infof("Downloading video from %s", youtubeURL)
			ys.logger.Infof("Saving from %s to %s ", youtubeURL, absFileName)
			err := stream.Download(youtubeURL, absFileName, ys.s.Config().DownloadLimit)
			if err != nil {
				ys.logger.Errorf("Got error while downloading video %s  ", err)
			}
			ys.logger.Infof("File %s saved for video %s", absFileName, youtubeURL)
		}
	}
	ys.s.waitGroup.Done()
}

// AddStream adds stream and runs it gracefully
func (ys *YoutubeStreamService) AddStream(stream data.StreamItem, download bool) {
	ys.logger.Infof("Adding a new stream: %s", stream.Name)
	ys.ss.Lock()
	ys.ss.Items[stream.ID] = stream
	ys.ss.Unlock()
	ys.logger.Infof("Starting streaming of %s", stream.Name)
	ys.s.waitGroup.Add(1)

	if !ys.s.Config().DisableStreaming {
		go ys.runStream(stream)
	}
	ys.s.waitGroup.Add(1)
	if download {
		go ys.downloadStream(stream)
	}
}

// UpdateStream updates stream storage
func (ys *YoutubeStreamService) UpdateStream(stream data.Stream, download bool) {
	ys.ss.Lock()
	item, ok := ys.ss.Items[stream.ID]
	if !ok {
		ys.logger.Errorf("%s does not exist in storage", stream.Name)
	}
	item.Lock()
	item.Links.Init()
	for _, link := range stream.Links {
		item.Links.PushBack(link.URL)
	}
	item.Unlock()

	item.Name = stream.Name
	ys.ss.Unlock()
	if download {
		go ys.downloadStream(stream.ToStreamItem())
	}
}

func (ys *YoutubeStreamService) runJobsForAutoStream(autoStream data.Stream, update bool) {

	streamData := ys.createStream(autoStream)
	streamData.IsAuto = true
	if len(streamData.Links) == 0 {
		return
	}
	if !update {
		ys.AddStream(streamData.ToStreamItem(), false)
	} else {
		ys.UpdateStream(streamData, false)
		ys.s.waitGroup.Add(1)
		go ys.runUpdateStream(autoStream)
	}
}

func (ys *YoutubeStreamService) createStream(as data.Stream) data.Stream {
	var links []data.StreamLink
	var streamData data.Stream
	streamData.ID = as.ID
	streamData.Name = as.Name
	streamData.Slug = as.Slug
	if as.Keywords != "" {
		for _, keyword := range strings.Split(as.Keywords, ",") {
			youtubeLinks, err := ys.getYoutubeContent(keyword, false, as.IsNews)
			if err != nil {
				ys.logger.Errorf("Error while requesting data, %v", err)
				continue
			}
			links = append(links, youtubeLinks...)
		}
	}
	if as.Channels != "" {

		for _, channel := range strings.Split(as.Channels, ",") {
			youtubeLinks, err := ys.getYoutubeContent(channel, true, as.IsNews)
			if err != nil {
				ys.logger.Errorf("Error while requesting data, %v", err)
				continue
			}
			links = append(links, youtubeLinks...)
		}
	}
	streamData.Links = ys.shuffleLinks(links)
	ys.logger.Info("Exiting")
	return streamData
}

func (ys *YoutubeStreamService) shuffleLinks(links []data.StreamLink) []data.StreamLink {
	for i := range links {
		j := rand.Intn(i + 1)
		links[i], links[j] = links[j], links[i]
	}
	return links
}

func (ys *YoutubeStreamService) getYoutubeContent(keyword string, isChannel, isNews bool) ([]data.StreamLink, error) {
	var links []data.StreamLink
	var results []*youtube.SearchResult
	var err error
	if isChannel {
		if isNews {
			results, err = ys.yc.SearchOnChannelByTime(keyword)
		} else {
			results, err = ys.yc.SearchOnChannel(keyword)
		}
	} else {
		results, err = ys.yc.Search(keyword)
	}
	if err != nil {
		return nil, err
	}
	for _, item := range results {
		switch item.Id.Kind {
		case "youtube#video":
			links = append(links, data.StreamLink{
				URL: fmt.Sprintf("https://youtube.com/watch?v=%s", item.Id.VideoId),
			})
		default:
			continue
		}
	}
	return links, nil
}

func (ys *YoutubeStreamService) runUpdateStream(as data.Stream) {

	for range time.Tick(time.Duration(as.UpdateFrequency) * time.Second) {
		ys.logger.Infof("Updating stream %s", as.Name)
		streamData := ys.createStream(as)

		ys.UpdateStream(streamData, false)
	}
	ys.s.waitGroup.Done()
}

func (ys *YoutubeStreamService) AddAutoStream(as data.Stream) {
	ys.runJobsForAutoStream(as, false)
}

func (ys *YoutubeStreamService) UpdateAutoStream(as data.Stream) {
	ys.runJobsForAutoStream(as, true)
}
