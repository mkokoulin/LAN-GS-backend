package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	s *http.Server
}

func New(handler *chi.Mux, addr string) *Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", addr),
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
				if err != nil {
					return nil, err
				}
				return &cert, nil
			},
		},
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
func (s *Server) StartTLS() error {
	err := s.s.ListenAndServeTLS("", "")
	if err != nil {
		return err
	}

	return nil
}
