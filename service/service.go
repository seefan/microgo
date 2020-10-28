package service

import (
	"reflect"
	"strings"
)

//Service base service
type Service interface {
}

//ServiceVersion base service
type ServiceVersion interface {
	//service version
	Version() string
}

//ServicePath base service
type ServicePath interface {
	//Based on this url
	Path() string
}

//ServiceName base service
type ServiceName interface {
	//Based on this url
	Name() string
}

func NewService(obj Service) *ServiceGroup {
	sg := &ServiceGroup{svr: obj}
	if s, ok := obj.(ServicePath); ok {
		sg.path = s.Path()
	} else {
		svr := reflect.TypeOf(obj)
		sg.path = strings.ToLower(svr.Elem().Name())
	}
	if s, ok := obj.(ServiceVersion); ok {
		sg.version = s.Version()
	}
	return sg
}

//ServiceGroup 服务组合
type ServiceGroup struct {
	path    string
	version string
	name    string
	svr     Service
}

//Path service path
func (s *ServiceGroup) Path() string {
	return s.path
}

//Version service version
func (s *ServiceGroup) Version() string {
	return s.version
}

//Service service
func (s *ServiceGroup) Service() Service {
	return s.svr
}

//Name service name
func (s *ServiceGroup) Name() Service {
	return s.svr
}
