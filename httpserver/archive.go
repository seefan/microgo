/*
@Time : 2019-01-19 10:50
@Author : seefan
@File : archive.go
@Software: microgo
*/
package httpserver

import (
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/service"
)

type archive struct {
	defaultUnit    *unit
	defaultVersion string
	svc            map[string]*unit
	begin          map[string][]service.Ware
	end            map[string][]service.Ware
}

func NewArchive() *archive {
	return &archive{
		svc:   make(map[string]*unit),
		begin: make(map[string][]service.Ware),
		end:   make(map[string][]service.Ware),
	}
}

func (a *archive) Put(sv service.Service) {
	a.svc[sv.Name()] = NewUnit(sv)
	if a.defaultVersion < sv.Version() {
		a.defaultVersion = sv.Version()
		a.defaultUnit = a.svc[sv.Name()]
	}
}
func (a *archive) BeginWare(name string, mid ...service.Ware) {
	if ms, ok := a.begin[name]; ok {
		a.begin[name] = append(ms, mid...)
	} else {
		a.begin[name] = mid
	}
}
func (a *archive) EndWare(name string, mid ...service.Ware) {
	if ms, ok := a.begin[name]; ok {
		a.end[name] = append(ms, mid...)
	} else {
		a.end[name] = mid
	}
}
func (a *archive) Get(v string) *unit {
	if sv, ok := a.svc[v]; ok {
		return sv
	} else {
		return a.defaultUnit
	}
}

func (a *archive) RunWare(name string, c ctx.Entry, wm map[string][]service.Ware) (err error) {
	if ms, ok := wm[name]; ok {
		for _, m := range ms {
			if m.Next != nil {
				err = m.Next(c)
			}
			if err != nil {
				return
			}
		}
	}
	return
}
