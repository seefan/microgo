package httpserver

import (
	"errors"
	"github.com/seefan/microgo/service"
	"log"
	"reflect"
	"runtime"
	"strings"
)

// Unit
type unit struct {
	svr     reflect.Value
	method  map[string]func(entry service.Entry) interface{}
	Version string
	Name    string
}

func (a *unit) resolve(s service.Service) {
	a.Version = s.Version()
	a.Name = s.Name()
	a.svr = reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		m := a.svr.MethodByName(t.Method(i).Name)
		if f, ok := m.Interface().(func(service.Entry) interface{}); ok {
			a.method[strings.ToLower(t.Method(i).Name)] = f
		}
	}
}
func NewUnit(s service.Service) *unit {
	a := &unit{
		method: make(map[string]func(entry service.Entry) interface{}),
	}
	a.resolve(s)
	return a
}

// RunMethod run a method
func (a *unit) RunMethod(name string, entry service.Entry) (re interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("RuntimeError")
			for i := 0; i < 10; i++ {
				funcName, file, line, ok := runtime.Caller(i)
				if ok {
					log.Printf("[func:%v,file:%v,line:%v]\n", runtime.FuncForPC(funcName).Name(), file, line)
				}
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
