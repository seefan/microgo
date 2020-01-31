package httpserver

import (
	"github.com/seefan/microgo/ctx"
	"net/http"
	"strings"
)

type archiveHandler struct {
	arch          *Archive
	call          func(interface{}, error, http.ResponseWriter)
	createContext func(httpContext *HTTPContext) ctx.Entry
}

//ServeHTTP server http method
func (a *archiveHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.ToLower(request.Method) == "options" {
		writer.WriteHeader(204)
		return
	}

	var result interface{}
	var err error
	defer func() {
		a.call(result, err, writer)
	}()
	var method, version string
	ms := strings.ToLower(request.URL.Path[a.arch.skip:])
	idx := strings.Index(ms, "/")
	if idx == -1 {
		method = ms
		version = a.arch.defaultMethodVersion[method]
	} else {
		method = ms[:idx]
		version = a.arch.getVersion(method, ms[idx+1:])
	}

	nc := newContext(writer, request)
	c := a.createContext(nc)
	if err = runWare(version, c, a.arch.before); err != nil {
		return
	}
	if result, err = a.arch.runMethod(method, version, c); err != nil {
		return
	}
	if err = runWare(version, c, a.arch.after); err != nil {
		return
	}
}
