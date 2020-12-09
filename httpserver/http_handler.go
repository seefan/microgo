package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/seefan/microgo/ctx"
)

type archiveHandler struct {
	arch          *Archive
	call          func(*ctx.Result, error, *http.Request, http.ResponseWriter)
	createContext func(httpContext *HTTPContext) ctx.Entry
}

//ServeHTTP server http method
func (a *archiveHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.ToLower(request.Method) == "options" {
		writer.WriteHeader(204)
		return
	}

	result := &ctx.Result{BeginNano: time.Now().UnixNano()}
	var err error

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
		result.EndNano = time.Now().UnixNano()
		a.call(result, err, request, writer)
		return
	}
	re, err := a.arch.runMethod(method, version, c)
	if err != nil {
		result.EndNano = time.Now().UnixNano()
		a.call(result, err, request, writer)
		return
	}
	result.Data = re
	if err = runWare(version, c, a.arch.after); err != nil {
		result.EndNano = time.Now().UnixNano()
		a.call(result, err, request, writer)
		return
	}
	result.EndNano = time.Now().UnixNano()
	a.call(result, err, request, writer)
}

type archiveWebsocketHandler struct {
	arch          *Archive
	call          func(*ctx.Result, error) []byte
	createContext func(httpContext *HTTPContext) ctx.Entry
}

//ServeHTTP server http method
func (a *archiveWebsocketHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(request, writer)
	if err != nil {
		a.call(nil, err)
		return
	}

	go func() {
		defer conn.Close()
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				a.call(nil, err)
				break
			}
			if op == ws.OpPing {
				_ = wsutil.WriteServerMessage(conn, ws.OpPong, []byte{})
				continue
			}
			r, err := a.handler(msg, writer, request)
			if err != nil {
				a.call(nil, err)
				break
			}
			bs := a.call(r, nil)
			err = wsutil.WriteServerMessage(conn, op, bs)
			if err != nil {
				a.call(nil, err)
				break
			}
		}
	}()

}
func (a *archiveWebsocketHandler) handler(msg []byte, writer http.ResponseWriter, request *http.Request) (result *ctx.Result, err error) {
	var method, version string
	result = &ctx.Result{BeginNano: time.Now().UnixNano()}
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
	nc.Set("_message", string(msg))
	c := a.createContext(nc)

	if err = runWare(version, c, a.arch.before); err != nil {
		result.EndNano = time.Now().UnixNano()
		return
	}
	re, err := a.arch.runMethod(method, version, c)
	if err != nil {
		result.EndNano = time.Now().UnixNano()
		return
	}
	result.Data = re
	err = runWare(version, c, a.arch.after)
	return
}
