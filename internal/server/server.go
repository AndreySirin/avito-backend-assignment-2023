package server

import (
	"fmt"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	lg         *logger.MyLogger
	httpServer *http.Server
}

func NewServer(lg *logger.MyLogger, adr string, hndl *HNDL) *Server {
	s := &Server{
		lg: lg,
	}
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/UserAddSegment", hndl.UserAddSegment)
			r.Post("/DeleteSegment", hndl.UserDeleteSegment)
		})
	})

	s.httpServer = &http.Server{
		Addr:    adr,
		Handler: r,
	}
	return s
}

func (s *Server) Start() error {
	s.lg.Info(fmt.Sprintf("Starting server on %s", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}
