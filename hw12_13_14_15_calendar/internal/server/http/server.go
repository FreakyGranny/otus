package internalhttp

import (
	"context"
	"io"
	"net/http"
)

// Server http server.
type Server struct {
	srv *http.Server
}

// Application business logic.
type Application interface {
	// TODO
}

// NewServer returns http server.
func NewServer(addr string, app Application) *Server {
	mux := http.NewServeMux()
	mux.Handle("/healthcheck", loggingMiddleware(http.HandlerFunc(HealthCheck)))

	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
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

// HealthCheck simple route.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK") //nolint
}
