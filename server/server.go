package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	log        *zap.Logger
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler, log *zap.Logger) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	s.log = log.Named("server")
	s.log.Info("Server started")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
