package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	s *http.Server
}

func New(addr string, handler *chi.Mux) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &Server {
		s: srv,
	}
}

func (s *Server) Start() (func(ctx context.Context) error, error) {
	err := s.s.ListenAndServe()
	if err != nil {
		return nil, err
	}

	return s.s.Shutdown, nil
}