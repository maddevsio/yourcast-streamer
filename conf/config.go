package conf

// StreamerConfig stores service configuration
type StreamerConfig struct {
	RTMPRootServerURL string
	WebUIURL          string
	FFMpegPath        string
	HTTPBindAddr      string
	RootPath          string
	YoutubeAPIKey     string
	DisableStreaming  bool
	DownloadLimit     int
}
