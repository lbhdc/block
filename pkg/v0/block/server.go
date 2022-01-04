package block

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Addr             string
	Port             uint32
	cancel           context.CancelFunc
	config           Configuration
	ctx              context.Context
	handlers         map[string]*HandlerPlugin
	mux              *http.ServeMux
	network          string
	s                *http.Server
	shutdownDuration time.Duration
}

func NewServer(ctx context.Context, options ...Option) *Server {
	s := new(Server)
	s.ctx, s.cancel = context.WithCancel(ctx)
	for _, o := range options {
		o(s)
	}
	if s.Port == 0 {
		s.Port = 8080
	}
	if s.handlers == nil {
		s.handlers = make(map[string]*HandlerPlugin)
	}
	if s.mux == nil {
		s.mux = http.NewServeMux()
	}
	if s.network == "" {
		s.network = "tcp"
	}
	if s.s == nil {
		s.s = new(http.Server)
	}
	if s.shutdownDuration == 0 {
		s.shutdownDuration = 30 * time.Second
	}
	return s
}

// Configure updates Server based on the config.toml settings used. Handlers described in the toml
// file are registered here.
func (s *Server) Configure(c *Configuration) {
	s.config = *c
	for _, handler := range c.Handlers {
		s.HandlePlugin(handler.Path, handler)
	}
}

// HandlePlugin registers a plugin handler with the server
func (s *Server) HandlePlugin(path string, cfg HandlerConfig) {
	s.handlers[path] = NewHandlerPlugin(s.ctx, cfg)
	s.mux.Handle(path, s.handlers[path])
}

// ListenAndServe starts all plugins before starting the block Server
func (s *Server) ListenAndServe() error {
	if err := s.StartPlugins(); err != nil {
		return err
	}
	lis, err := net.Listen(s.network, fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return err
	}
	defer func(lis net.Listener) {
		if err = lis.Close(); err != nil {
			log.Fatalln(err)
		}
	}(lis)
	return s.Serve(lis)
}

func (s *Server) Serve(lis net.Listener) error {
	s.s.Handler = s.mux
	return s.s.Serve(lis)
}

// Shutdown gracefully shutsdown the http server before killing all child processes
func (s *Server) Shutdown() error {
	s.cancel() // kill the server context to allow everything to stop
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownDuration)
	defer cancel()
	// stop new requests from coming in and shut down the external server
	if err := s.s.Shutdown(ctx); err != nil {
		return err
	}
	// kill each sub process
	for path := range s.handlers {
		if err := s.handlers[path].Kill(); err != nil {
			log.Println(err)
		}
	}
	return nil
}

// StartPlugins launches all plugins as child processes
func (s *Server) StartPlugins() error {
	for path := range s.handlers {
		if err := s.handlers[path].Start(s.config.ConfigPath); err != nil {
			log.WithError(err).Error("s.Handlers[path].Start")
			return err
		}
	}
	return nil
}
