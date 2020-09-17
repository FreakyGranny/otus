package internalhttp

import (
	"context"
	"net/http"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
)

// Server http server.
type Server struct {
	srv *http.Server
}

// NewServer returns http server.
func NewServer(addr string, app app.Application) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: loggingMiddleware(NewEventHandler(app)),
		},
	}
}

// Start starts http server.
func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop stops http server.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
