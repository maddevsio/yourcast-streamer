package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/gen1us2k/log"
	"github.com/maddevsio/yourcast-streamer/conf"
	"github.com/maddevsio/yourcast-streamer/service"
	"github.com/urfave/cli"
)

// LogLevel used for logger configuration. For more detailed information see https://github.com/gen1us2k/log
// RTMPRootServerURL used for configuration of url where to stream from youtube
// WebUIURL used for configuration of http web api
// FFMpegPath used to store ffmpeg path to binary
// HTTPBindAddr used for configuration of inner HTTP api where to bind to
// RootPath used for configuration where to store files, downloaded via ffmpeg while streaming
var (
	LogLevel          string
	RTMPRootServerURL string
	WebUIURL          string
	FFMpegPath        string
	HTTPBindAddr      string
	RootPath          string
	YoutubeAPIKey     string
	DownloadLimit     int
	DisableStreaming  bool
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1dev"
	app.Name = "Youtube-Reastreamer"
	app.Usage = "Restreams video directly from youtube"
	app.Action = actionStart
	app.Flags = []cli.Flag{

		cli.IntFlag{
			Name:        "download_limit",
			Value:       100,
			EnvVar:      "RESTREAMER_DOWNLOAD_LIMIT",
			Destination: &DownloadLimit,
		},
		cli.StringFlag{
			Name:        "rtmp_server_url",
			Value:       "rtmp://localhost/hls",
			EnvVar:      "RESTREAMER_RTMP_ROOT_SERVER_URL",
			Destination: &RTMPRootServerURL,
		},
		cli.StringFlag{
			Name:        "log_level",
			Value:       "debug",
			EnvVar:      "RESTREAMER_LOG_LEVEL",
			Destination: &LogLevel,
		},
		cli.StringFlag{
			Name:        "web_ui_url",
			Value:       "http://localhost:8000",
			EnvVar:      "RESTREAMER_WEB_UI_URL",
			Destination: &WebUIURL,
		},
		cli.StringFlag{
			Name:        "ffmpeg_path",
			Value:       "/usr/local/bin/ffmpeg",
			EnvVar:      "RESTREAMER_FFMPEG_PATH",
			Destination: &FFMpegPath,
		},
		cli.StringFlag{
			Name:        "http_bind_addr",
			Value:       ":8080",
			EnvVar:      "RESTREAMER_HTTP_BIND_ADDR",
			Destination: &HTTPBindAddr,
		},
		cli.StringFlag{
			Name:        "youtube_api_key",
			Value:       "Aiza...",
			EnvVar:      "RESTREAMER_YOUTUBE_API_KEY",
			Destination: &YoutubeAPIKey,
		},
		cli.StringFlag{
			Name:        "root_path",
			Value:       "./storage",
			EnvVar:      "RESTREAMER_FILE_ROOT_PATH",
			Destination: &RootPath,
		},
		cli.BoolFlag{
			Name:        "disable_streaming",
			EnvVar:      "RESTREAMER_DISABLE_STREAMING",
			Destination: &DisableStreaming,
		},
	}
	app.Before = func(ctx *cli.Context) error {
		log.SetLevel(log.MustParseLevel(LogLevel))
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionStart(ctx *cli.Context) error {
	log.Info("Checking storage directory")

	if _, err := os.Stat(RootPath); os.IsNotExist(err) {
		os.Mkdir(RootPath, 0755)
	}
	conf := &conf.StreamerConfig{
		RTMPRootServerURL: RTMPRootServerURL,
		WebUIURL:          WebUIURL,
		FFMpegPath:        FFMpegPath,
		HTTPBindAddr:      HTTPBindAddr,
		RootPath:          RootPath,
		YoutubeAPIKey:     YoutubeAPIKey,
		DisableStreaming:  DisableStreaming,
		DownloadLimit:     DownloadLimit,
	}
	log.Info("Starting streamer...")
	streamer := service.NewStreamer(conf)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	defer signal.Stop(signalChan)

	go func() {
		<-signalChan
		log.Info("signal received, stopping...")
		streamer.Stop()

		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	err := streamer.Start()

	if err != nil {
		log.Fatalf("error on local node start, %v", err)
	}

	streamer.WaitStop()
	return nil
}
