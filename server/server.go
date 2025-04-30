package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"myapp/stor"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logg       *slog.Logger
	Sub        stor.User_Subscription
}

func NewServer(logger *slog.Logger, adr string, subscription stor.User_Subscription) *Server {
	s := &Server{
		logg: logger,
		Sub:  subscription,
	}
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/UserAddSegment", s.UserAddSegment)
			r.Post("/DeleteSegment", s.UserDeleteSegment)
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
