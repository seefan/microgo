package httpserver

import (
	"reflect"
	"strings"

	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/service"
)

// Unit
type unit struct {
	method  map[string]Method
	Version string
	Path    string
}

func (a *unit) resolve(s *service.ServiceGroup) {
	a.Version = s.Version()
	if strings.HasPrefix(s.Path(), "/") {
		a.Path = s.Path()[1:]
	} else {
		a.Path = s.Path()
	}
	svr := reflect.ValueOf(s.Service())
	t := reflect.TypeOf(s.Service())
	for i := 0; i < t.NumMethod(); i++ {
		m := svr.MethodByName(t.Method(i).Name)
		if f, ok := m.Interface().(func(entry ctx.Entry) interface{}); ok {
			a.method[strings.ToLower(t.Method(i).Name)] = f
		}
		if f, ok := m.Interface().(func(entry ctx.Entry)); ok {
			a.method[strings.ToLower(t.Method(i).Name)] = f
		}
		if f, ok := m.Interface().(func() interface{}); ok {
			a.method[strings.ToLower(t.Method(i).Name)] = f
		}
	}
}
func newUnit(s *service.ServiceGroup) *unit {
	a := &unit{
		method: make(map[string]Method),
	}
	a.resolve(s)
	return a
}
