package bot

import (
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// YoutubeClient stores copy of youtube service
type YoutubeClient struct {
	youtubeService *youtube.Service
}

// NewYoutubeClient creates client for youtube
func NewYoutubeClient(devKey string) (*YoutubeClient, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: devKey},
	}
	yc := new(YoutubeClient)
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}
	yc.youtubeService = service
	return yc, nil
}

// Search performs query to youtube and return results
func (yc *YoutubeClient) Search(query string) ([]*youtube.SearchResult, error) {
	time := time.Now().AddDate(0, 0, -2).Format(time.RFC3339)
	call := yc.youtubeService.Search.List("id,snippet").
		Q(query).
		MaxResults(50).
		PublishedAfter(time).
		Order("date")
	response, err := call.Do()
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

// SearchOnChannel performs query to youtube and return results
func (yc *YoutubeClient) SearchOnChannel(query string) ([]*youtube.SearchResult, error) {
	channelCall := yc.youtubeService.Channels.List("id").ForUsername(query)
	response, err := channelCall.Do()
	if err != nil {
		return nil, err
	}
	searchQuery := query
	if len(response.Items) > 0 {
		searchQuery = response.Items[0].Id
	}
	call := yc.youtubeService.Search.List("id,snippet").
		ChannelId(searchQuery).
		MaxResults(50).
		Order("date")
	ChannelResponse, err := call.Do()
	if err != nil {
		return nil, err
	}
	return ChannelResponse.Items, nil
}

// SearchOnChannelByTime performs query to youtube and return results
func (yc *YoutubeClient) SearchOnChannelByTime(query string) ([]*youtube.SearchResult, error) {
	channelCall := yc.youtubeService.Channels.List("id").ForUsername(query)
	response, err := channelCall.Do()
	time := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	if err != nil {
		return nil, err
	}
	searchQuery := query
	if len(response.Items) > 0 {
		searchQuery = response.Items[0].Id
	}
	call := yc.youtubeService.Search.List("id,snippet").
		ChannelId(searchQuery).
		MaxResults(50).
		PublishedAfter(time).
		Order("date")
	ChannelResponse, err := call.Do()
	if err != nil {
		return nil, err
	}
	return ChannelResponse.Items, nil
}
