package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logg       *slog.Logger
}

func NewServer(logger *slog.Logger, adr string) *Server {
	s := &Server{
		logg: logger,
	}
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
		})
	})

	s.httpServer = &http.Server{
		Addr:    adr,
		Handler: r,
	}
	return s
}

func (s *Server) Start() error {
	s.logg.Info(fmt.Sprintf("Starting server on %s", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}
