package service

import (
	"fmt"
	"sync"

	"github.com/gen1us2k/log"
	"github.com/maddevsio/yourcast-streamer/conf"
)

// Streamer is main struct of daemon
// it stores all services that used by
type Streamer struct {
	config *conf.StreamerConfig

	services  map[string]Service
	waitGroup sync.WaitGroup

	logger log.Logger
}

// NewStreamer creates and returns new StreamerInstance
func NewStreamer(config *conf.StreamerConfig) *Streamer {
	s := new(Streamer)
	s.config = config
	s.logger = log.NewLogger("streamer_worker")
	s.services = make(map[string]Service)
	s.AddService(&YoutubeStreamService{})
	s.AddService(&HTTPService{})
	return s
}

// Start starts all services in separate goroutine
func (s *Streamer) Start() error {
	s.logger.Info("Starting streamer backend")
	for _, service := range s.services {
		s.logger.Infof("Initializing: %s\n", service.Name())
		if err := service.Init(s); err != nil {
			return fmt.Errorf("initialization of %q finished with error: %v", service.Name(), err)
		}
		s.waitGroup.Add(1)

		go func(srv Service) {
			defer s.waitGroup.Done()
			s.logger.Infof("running %q service\n", srv.Name())
			if err := srv.Run(); err != nil {
				s.logger.Errorf("error on run %q service, %v", srv.Name(), err)
			}
		}(service)
	}
	return nil
}

// AddService adds service into Streamer.services map
func (s *Streamer) AddService(srv Service) {
	s.services[srv.Name()] = srv

}

// Config returns current instance of StreamerConfig
func (s *Streamer) Config() conf.StreamerConfig {
	return *s.config
}

// Stop stops all services running
func (s *Streamer) Stop() {
	s.logger.Info("Worker is stopping...")
	for _, service := range s.services {
		service.Stop()
	}
}

// WaitStop blocks main thread and waits when all goroutines will be stopped
func (s *Streamer) WaitStop() {
	s.waitGroup.Wait()
}

// YoutubeStreamService returns *YoutubeStreamService
func (s *Streamer) YoutubeStreamService() *YoutubeStreamService {
	service, ok := s.services["youtube_stream_service"]
	if !ok {
		s.logger.Info("youtube_stream_service not found")
	}
	return service.(*YoutubeStreamService)
}
