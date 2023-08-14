package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	s *http.Server
}

func New(handler *chi.Mux) *Server {
	srv := &http.Server{
		Addr:    ":8000",
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

// StartTLS is the method to start the server with tls
func (s *Server) StartTLS(certFile, keyFile string) error {
	err := s.s.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		return err
	}

	return nil
}
