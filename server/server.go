package server

import (
	"context"
	"fmt"
)

type Runnable interface {
	Start(ctx context.Context) error
	Stop() error
}

// Server for basic
type Server struct {
	Host     string
	Port     int
	Name     string
	InitFunc func()
}

// Init init server
func (s *Server) Init(host string, port int, name ...string) {
	if len(name) > 0 {
		s.Name = name[0]
	}
	s.Host = host
	s.Port = port
}

// Address get server address
func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// String return a short paragraph of explanatory text
func (s *Server) String() string {
	return fmt.Sprintf("%s at %s:%d", s.Name, s.Host, s.Port)
}
