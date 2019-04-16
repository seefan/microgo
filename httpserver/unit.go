package httpserver

import (
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/service"
	"reflect"
	"strings"
)

// Unit
type unit struct {
	method  map[string]func(entry ctx.Entry) interface{}
	Version string
	Path    string
}

func (a *unit) resolve(s service.Service) {
	a.Version = s.Version()
	if strings.HasPrefix(s.Path(), "/") {
		a.Path = s.Path()[1:]
	} else {
		a.Path = s.Path()
	}
	svr := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		m := svr.MethodByName(t.Method(i).Name)
		if f, ok := m.Interface().(func(ctx.Entry) interface{}); ok {
			a.method[strings.ToLower(t.Method(i).Name)] = f
		}
	}
}
func newUnit(s service.Service) *unit {
	a := &unit{
		method: make(map[string]func(entry ctx.Entry) interface{}),
	}
	a.resolve(s)
	return a
}
