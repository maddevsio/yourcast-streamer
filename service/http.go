package service

import (
	"encoding/json"

	"github.com/gen1us2k/log"
	"github.com/labstack/echo"
	"github.com/maddevsio/yourcast-streamer/service/data"
)

// HTTPService implements a simple api for streamer
type HTTPService struct {
	BaseService

	s      *Streamer
	e      *echo.Echo
	ys     *YoutubeStreamService
	logger log.Logger
}

// Name returns name of service
func (h *HTTPService) Name() string {
	return "http_api"
}

// Init initializes echo server, http routing and logger
func (h *HTTPService) Init(s *Streamer) error {
	h.s = s
	h.logger = log.NewLogger(h.Name())
	h.e = echo.New()
	h.ys = h.s.YoutubeStreamService()
	h.e.POST("/stream/add", h.addStream)
	h.e.POST("/stream/update", h.updateStream)
	return nil
}

// Run runs service
func (h *HTTPService) Run() error {
	h.e.Start(h.s.Config().HTTPBindAddr)
	return nil
}

func (h *HTTPService) addStream(c echo.Context) error {
	streamData := c.FormValue("data")
	var stream data.Stream
	err := json.Unmarshal([]byte(streamData), &stream)
	if err != nil {
		h.logger.Errorf("caught error on json unmarshaling: %s", err)
		return err
	}
	if !stream.IsAutoStream() {
		h.logger.Infof("adding a new stream: %s", stream.Name)
		h.ys.AddStream(stream.ToStreamItem(), true)
	} else {
		h.logger.Infof("adding autostream: %s", stream.Name)
		h.ys.AddAutoStream(stream)
	}
	return nil
}

func (h *HTTPService) updateStream(c echo.Context) error {
	streamData := c.FormValue("data")
	h.logger.Info(streamData)
	var stream data.Stream
	err := json.Unmarshal([]byte(streamData), &stream)
	if err != nil {
		h.logger.Errorf("caught error on json unmarshaling: %s", err)
		return err
	}
	if !stream.IsAutoStream() {
		h.logger.Infof("updating a stream: %s", stream.Name)
		h.ys.UpdateStream(stream, true)
	} else {

		h.logger.Infof("adding autostream: %s", stream.Name)
		h.ys.UpdateAutoStream(stream)
	}
	return nil
}
