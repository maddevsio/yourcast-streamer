package stream

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
)

// FromYoutubeURL streams from youtube.
// First we need to get streamable url to file, then restream it
func FromYoutubeURL(ffmpeg, youtubeURL, dst, rootPath string) error {
	_, url, err := GetStreamURL(youtubeURL)
	if err != nil {
		return fmt.Errorf("Got error %s while getting streamable  youtube url for video %s ", err, youtubeURL)
	}
	ffmpegArgs := []string{
		"-re", "-i", url,
		"-headers", "User-Agent: Go-http-client/1.1",
		"-c", "copy", "-f", "flv", dst,
	}

	cmd := exec.Command(ffmpeg, ffmpegArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// FromLocalFile streams video from local file downloaded from youtube
func FromLocalFile(ffmpeg, fileName, dst string) error {

	ffmpegArgs := []string{
		"-re", "-i", fileName,
		"-c", "copy", "-f", "flv", dst,
	}
	cmd := exec.Command(ffmpeg, ffmpegArgs...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// GetFileNameByURL returns MD5 hash from youtube url
func GetFileNameByURL(url string, rootPath string) string {
	fileName := fmt.Sprintf("%x.mp4", md5.Sum([]byte(url)))
	return fmt.Sprintf("%s/%s", rootPath, fileName)
}

// RemoveFile removes file by youtube URL
func RemoveFile(fileName string) error {
	return os.Remove(fmt.Sprintf("%s.download", fileName))
}

// RenameFile removes .download extension from filename
func RenameFile(fileName string) error {
	downloadedFile := fmt.Sprintf("%s.download", fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if _, err := os.Stat(downloadedFile); err == nil {
			return os.Rename(downloadedFile, fileName)
		}
	}
	return nil
}

// FileExist checks if file exist in current directory
func FileExist(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}
