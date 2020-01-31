/*
@Time : 2019-01-19 10:50
@Author : seefan
@File : archive.go
@Software: microgo
*/
package httpserver

import (
	"github.com/seefan/goerr"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/service"
	"runtime"
	"strings"
)

//Archive Archive
type Archive struct {
	defaultMethodVersion map[string]string
	currentVersion       string
	svc                  map[string]*unit          //version:service
	method               map[string]Method         //version:method
	before               map[string][]service.Ware //version:ware
	after                map[string][]service.Ware //version:ware
	skip                 int
	name                 string
}

func newArchive(name string) *Archive {
	return &Archive{
		name:                 name,
		skip:                 len(name),
		svc:                  make(map[string]*unit),
		defaultMethodVersion: make(map[string]string),
		method:               make(map[string]Method),
		before:               make(map[string][]service.Ware),
		after:                make(map[string][]service.Ware),
	}
}

//Put set service and version
func (a *Archive) Put(sv service.Service) {
	t := newUnit(sv)
	a.svc[sv.Path()] = t
	a.currentVersion = strings.ToLower(sv.Version())
	for name, f := range t.method {
		name = strings.ToLower(name)
		if _, ok := a.method[name]; !ok {
			a.method[name] = make(Method)
		}
		a.method[name][a.currentVersion] = f
		if v, ok := a.defaultMethodVersion[name]; !ok || a.currentVersion > v {
			a.defaultMethodVersion[name] = a.currentVersion
		}
	}
}

//Before set before ware
func (a *Archive) Before(mid service.Ware, svc ...service.Service) {
	var cv string
	if len(svc) > 0 {
		cv = strings.ToLower(svc[0].Version())
	} else {
		cv = a.currentVersion
	}
	if ms, ok := a.before[cv]; ok {
		a.before[cv] = append(ms, mid)
	} else {
		a.before[cv] = []service.Ware{mid}
	}
}

//After set after ware
func (a *Archive) After(mid service.Ware, svc ...service.Service) {
	var cv string
	if len(svc) > 0 {
		cv = strings.ToLower(svc[0].Version())
	} else {
		cv = a.currentVersion
	}
	if ms, ok := a.before[cv]; ok {
		a.after[cv] = append(ms, mid)
	} else {
		a.after[cv] = []service.Ware{mid}
	}
}

func (a *Archive) getVersion(name, version string) string {
	if mv, ok := a.method[name]; ok {
		if _, ok := mv[version]; ok {
			return version
		}
		return a.defaultMethodVersion[name]
	}
	return ""
}

// RunMethod run a method
func (a *Archive) runMethod(name, version string, entry ctx.Entry) (re interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			ne := goerr.String("RuntimeError:%s", e)
			get := false
			for i := 0; i < 10; i++ {
				if fp, f1, l, ok := runtime.Caller(i); ok {
					if get {
						ne.AttachE(goerr.String(runtime.FuncForPC(fp).Name()).File(f1).Line(l))
					}
					if strings.Index(f1, "github.com/seefan/microgo/httpserver/archive.go") != -1 {
						get = true
					}
				}
			}
			err = ne
		}
	}()
	if name == "" { //check default method
		if mv, ok := a.method["default"]; ok {
			return mv[a.defaultMethodVersion["default"]](entry), nil
		}
	}
	if mv, ok := a.method[name]; ok {
		if m, ok := mv[version]; ok {
			return m(entry), nil
		} else {
			return mv[a.defaultMethodVersion[name]](entry), nil
		}
	}

	return nil, goerr.Errorf(goerr.String("Methods:%s Version:%s", name, version), "MethodNotFound")
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
