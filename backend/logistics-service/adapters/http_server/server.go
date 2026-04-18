package http_server

import (
	"context"
	"errors"
	"net/http"
	"time"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/ports"
)

type Server struct {
	Logger ports.Logger
	server *http.Server
}

func New(address string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	return &Server{server: &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return coreErrors.ErrListenAndServeServer
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return coreErrors.ErrShutdownServer
	}
	return nil
}