package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Just-maple/godi"
	"github.com/Just-maple/xmux/examples/webapp/pkg/app"
	"github.com/Just-maple/xmux/examples/webapp/pkg/controller"
	"github.com/Just-maple/xmux/examples/webapp/pkg/di"
)

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	config     ServerConfig
	httpServer *http.Server
	container  *godi.Container
	shutdown   func(context.Context, bool)
}

func NewServer(config ServerConfig) (*Server, error) {
	container, shutdown, err := di.BuildContainer()
	if err != nil {
		return nil, err
	}

	return &Server{
		config:    config,
		container: container,
		shutdown:  shutdown,
	}, nil
}

func (s *Server) Start() error {
	app := app.NewApplication(s.container)

	ctrl := controller.NewController()
	app.RegisterRoutes(ctrl)

	s.httpServer = &http.Server{
		Addr:         s.config.Addr,
		Handler:      ctrl,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	log.Printf("Server starting on %s", s.config.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}

	if s.shutdown != nil {
		s.shutdown(ctx, false)
	}

	log.Println("Server stopped")
	return nil
}

func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.Shutdown(shutdownCtx)
	}

	return nil
}
