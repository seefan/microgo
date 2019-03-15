package httpserver

import (
	"errors"
	"github.com/seefan/goerr"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/service"
	"reflect"
	"runtime"
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
func NewUnit(s service.Service) *unit {
	a := &unit{
		method: make(map[string]func(entry ctx.Entry) interface{}),
	}
	a.resolve(s)
	return a
}

// RunMethod run a method
func (a *unit) RunMethod(name string, entry ctx.Entry) (re interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			ne := goerr.String("RuntimeError")
			_, file, line, ok := runtime.Caller(3)
			if ok {
				err = ne.Line(line).File(file)
			}
		}
	}()
	m, ok := a.method[strings.ToLower(name)]
	if !ok {
		err = errors.New("MethodNotFound")
		return
	}
	re = m(entry)
	return
}
