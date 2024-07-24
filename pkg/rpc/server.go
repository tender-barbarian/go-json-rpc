package rpc

import (
	"net/http"
)

type Server struct {
	URL     string
	Handler *Handler
}

func New(URL string) (*Server, error) {
	return &Server{URL: URL, Handler: &Handler{}}, nil
}

func (s *Server) Start() error {
	http.Handle("/", s.Handler)

	err := http.ListenAndServe(s.URL, nil)
	if err != nil {
		return err
	}

	return nil
}
