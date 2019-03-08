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
	svc            map[string]*unit          //version:service
	begin          map[string][]service.Ware //version:ware
	end            map[string][]service.Ware //version:ware
}

func NewArchive() *archive {
	return &archive{
		svc:   make(map[string]*unit),
		begin: make(map[string][]service.Ware),
		end:   make(map[string][]service.Ware),
	}
}

func (a *archive) Put(sv service.Service) {
	a.svc[sv.Path()] = NewUnit(sv)
	if a.defaultVersion < sv.Version() {
		a.defaultVersion = sv.Version()
		a.defaultUnit = a.svc[sv.Path()]
	}
}
func (a *archive) BeginWare(svc service.Service, mid ...service.Ware) {
	if ms, ok := a.begin[svc.Version()]; ok {
		a.begin[svc.Version()] = append(ms, mid...)
	} else {
		a.begin[svc.Version()] = mid
	}
}
func (a *archive) EndWare(svc service.Service, mid ...service.Ware) {
	if ms, ok := a.begin[svc.Version()]; ok {
		a.end[svc.Version()] = append(ms, mid...)
	} else {
		a.end[svc.Version()] = mid
	}
}
func (a *archive) Get(v string) *unit {
	if sv, ok := a.svc[v]; ok {
		return sv
	} else {
		return a.defaultUnit
	}
}

func runWare(version string, c ctx.Entry, wm map[string][]service.Ware) (err error) {
	if ms, ok := wm[version]; ok {
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
