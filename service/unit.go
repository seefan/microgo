package service

import (
	"errors"
	"log"
	"reflect"
	"runtime"
	"strings"
)

// Unit
type Unit struct {
	svr     reflect.Value
	method  map[string]reflect.Value
	Version string
	Name    string
}

func (a *Unit) resolve(s Service) {
	a.Version = s.Version()
	a.Name = s.Name()
	a.svr = reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		a.method[strings.ToLower(t.Method(i).Name)] = a.svr.MethodByName(t.Method(i).Name)
	}
}
func newUnit(s Service) *Unit {
	a := &Unit{
		method: make(map[string]reflect.Value),
	}
	a.resolve(s)
	return a
}

// RunMethod run a method
func (a *Unit) RunMethod(name string, param ...interface{}) (re []interface{}, err error) {
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
	var in []reflect.Value
	for _, v := range param {
		in = append(in, reflect.ValueOf(v))
	}
	result := m.Call(in)
	for _, r := range result {
		re = append(re, r.Interface())
	}
	return
}
