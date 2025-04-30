package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/storage"
)

type Server struct {
	lg         *slog.Logger
	httpServer *http.Server
	Sub        storage.User_Subscription
}

func NewServer(logger *slog.Logger, adr string, subscription storage.User_Subscription) *Server {
	s := &Server{
		lg:  logger,
		Sub: subscription,
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
