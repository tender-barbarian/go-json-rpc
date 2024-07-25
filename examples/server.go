package main

import (
	"net/http"

	"github.com/tender-barbarian/go-json-rpc/pkg/rpc"
)

type Server struct {
	URL     string
	Handler *rpc.Handler
}

func New(URL string) (*Server, error) {
	return &Server{URL: URL, Handler: &rpc.Handler{}}, nil
}

func (s *Server) Start() error {
	http.Handle("/", s.Handler)

	err := http.ListenAndServe(s.URL, nil)
	if err != nil {
		return err
	}

	return nil
}
