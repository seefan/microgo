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
	currentVersion string
	svc            map[string]*unit          //version:service
	before         map[string][]service.Ware //version:ware
	after          map[string][]service.Ware //version:ware
}

func NewArchive() *archive {
	return &archive{
		svc:    make(map[string]*unit),
		before: make(map[string][]service.Ware),
		after:  make(map[string][]service.Ware),
	}
}

// set service and version
func (a *archive) Put(sv service.Service) {
	a.svc[sv.Path()] = NewUnit(sv)
	a.currentVersion = sv.Version()
	if a.defaultVersion < sv.Version() {
		a.defaultVersion = sv.Version()
		a.defaultUnit = a.svc[sv.Path()]
	}
}

// set before ware
func (a *archive) Before(mid service.Ware, svc ...service.Service) {
	var cv string
	if len(svc) > 0 {
		cv = svc[0].Version()
	} else {
		cv = a.currentVersion
	}
	if ms, ok := a.before[cv]; ok {
		a.before[cv] = append(ms, mid)
	} else {
		a.before[cv] = []service.Ware{mid}
	}
}

// set after ware
func (a *archive) After(mid service.Ware, svc ...service.Service) {
	var cv string
	if len(svc) > 0 {
		cv = svc[0].Version()
	} else {
		cv = a.currentVersion
	}
	if ms, ok := a.before[cv]; ok {
		a.after[cv] = append(ms, mid)
	} else {
		a.after[cv] = []service.Ware{mid}
	}
}

// get service
func (a *archive) Get(v string) *unit {
	if sv, ok := a.svc[v]; ok {
		return sv
	} else {
		return a.defaultUnit
	}
}

// run ware
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
