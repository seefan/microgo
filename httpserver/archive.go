/*
@Time : 2019-01-19 10:50
@Author : seefan
@File : archive.go
@Software: microgo
*/
package httpserver

import "github.com/seefan/microgo/service"

type archive struct {
	defaultUnit    *unit
	defaultVersion string
	svc            map[string]*unit
}

func NewArchive() *archive {
	return &archive{
		svc: make(map[string]*unit),
	}
}

func (a *archive) Put(sv service.Service) {
	a.svc[sv.Name()] = NewUnit(sv)
	if a.defaultVersion < sv.Version() {
		a.defaultVersion = sv.Version()
		a.defaultUnit = a.svc[sv.Name()]
	}
}
func (a *archive) Get(v string) *unit {
	if sv, ok := a.svc[v]; ok {
		return sv
	} else {
		return a.defaultUnit
	}
}
