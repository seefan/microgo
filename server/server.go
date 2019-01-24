package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/seefan/microgo/service"
)

// Server for basic
type Server struct {
	Host       string
	Port       int
	Name       string
	serviceMap map[string]*service.Archive
	a          time.Time
}

// Init init server
func (s *Server) Init(host string, port int, name ...string) {
	s.serviceMap = make(map[string]*service.Archive)
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

// Register the service with the server
//
func (s *Server) Register(svc service.Service) {
	if _, ok := s.serviceMap[svc.Name()]; !ok {
		s.serviceMap[svc.Name()] = service.NewArchive()
	}
	s.serviceMap[svc.Name()].PutService(svc)
}

// GetServiceArchive get service archive
func (s *Server) GetService(name string) (*service.Archive, error) {
	if svc, ok := s.serviceMap[name]; ok {
		return svc, nil
	} else {
		return nil, errors.New("UnknownService")
	}
}

// String return a short paragraph of explanatory text
func (s *Server) String() string {
	return fmt.Sprintf("%s at %s:%d", s.Name, s.Host, s.Port)
}
