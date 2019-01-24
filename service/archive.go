/*
@Time : 2019-01-19 10:50
@Author : seefan
@File : archive.go
@Software: microgo
*/
package service

type Archive struct {
	defaultUnit    *Unit
	defaultVersion string
	svr            map[string]*Unit
}

func NewArchive() *Archive {
	return &Archive{
		svr: make(map[string]*Unit),
	}
}

func (a *Archive) PutService(sv Service) {
	s := newUnit(sv)
	a.svr[s.Version] = s
	if a.defaultVersion < s.Version {
		a.defaultVersion = s.Version
		a.defaultUnit = s
	}
}
func (a *Archive) GetUnit(v string) *Unit {
	if sv, ok := a.svr[v]; ok {
		return sv
	} else {
		return a.defaultUnit
	}
}
