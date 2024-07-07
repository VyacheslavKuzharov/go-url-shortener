package httpserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	defaultReadTimeout  = 40 * time.Second
	defaultWriteTimeout = 40 * time.Second
	DefaultHost         = "localhost"
	DefaultPort         = "8080"
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &Server{
		server: httpServer,
	}

	return s
}

func (s *Server) Start(host string, port string) {
	s.server.Addr = net.JoinHostPort(host, port)
	fmt.Println("Starting server on:", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
