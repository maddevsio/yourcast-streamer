package stream

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/otium/ytdl"

	"github.com/mxk/go-flowrate/flowrate"
)

// GetStreamURL returns Title of youtube video and url for re-stream
func GetStreamURL(url string) (string, string, error) {
	info, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return "", "", err
	}
	foundFormat := getBestFormat(info.Formats)
	videoURL, err := info.GetDownloadURL(foundFormat)
	if err != nil {
		return "", "", err
	}
	return info.Title, videoURL.String(), err
}

func getBestFormat(formats ytdl.FormatList) ytdl.Format {
	var foundFormat ytdl.Format

	bestFormats := formats.Filter(ytdl.FormatResolutionKey, []interface{}{"480p", "360p"}).Filter(ytdl.FormatExtensionKey, []interface{}{"mp4"}).Filter(ytdl.FormatAudioEncodingKey, []interface{}{"aac"}).Extremes(ytdl.FormatResolutionKey, true).Extremes(ytdl.FormatAudioBitrateKey, true)

	for _, format := range bestFormats {
		if format.Extension == "mp4" {
			foundFormat = format
		}
	}
	return foundFormat
}

func Download(youtubeURL, fileName string, downloadLimit int) error {
	dst := fmt.Sprintf("%s.download", fileName)
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := ytdl.GetVideoInfo(youtubeURL)
	if err != nil {
		os.Remove(dst)
		return err
	}
	foundFormat := getBestFormat(info.Formats)

	u, err := info.GetDownloadURL(foundFormat)
	if err != nil {
		os.Remove(dst)
		return err
	}
	resp, err := http.Get(u.String())
	if err != nil {
		os.Remove(dst)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}
	wrappedIn := flowrate.NewReader(resp.Body, int64(downloadLimit)*1024)
	_, err = io.Copy(file, wrappedIn)
	if err != nil {
		os.Remove(dst)
		return err
	}
	RenameFile(fileName)

	return nil
}

func contains(resolution string, resolutions []string) bool {
	for _, res := range resolutions {
		if res == resolution {
			return true
		}
	}
	return false
}
