package httpserver

import (
	"log"
	"net/http"
	"time"
)

const (
	defaultReadTimeout  = 40 * time.Second
	defaultWriteTimeout = 40 * time.Second
	defaultAddr         = ":8080"
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Server{
		server: httpServer,
	}

	return s
}

func (s *Server) Start() {
	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
