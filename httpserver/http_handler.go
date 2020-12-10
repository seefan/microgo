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

	content := &ctx.Result{BeginNano: time.Now().UnixNano()}
	var err error

	ms := strings.ToLower(request.URL.Path[a.arch.skip:])
	idx := strings.Index(ms, "/")
	if idx == -1 {
		content.Method = ms
		content.Version = a.arch.defaultMethodVersion[content.Method]
	} else {
		content.Method = ms[:idx]
		content.Version = a.arch.getVersion(content.Method, ms[idx+1:])
	}
	nc := newContext(writer, request)

	content.Request = a.createContext(nc)
	if err = runWare(content.Version, content.Request, a.arch.before); err != nil {
		content.EndNano = time.Now().UnixNano()
		a.call(content, err, request, writer)
		return
	}
	re, err := a.arch.runMethod(content.Method, content.Version, content.Request)
	if err != nil {
		content.EndNano = time.Now().UnixNano()
		a.call(content, err, request, writer)
		return
	}
	content.Response = re
	if err = runWare(content.Version, content.Request, a.arch.after); err != nil {
		content.EndNano = time.Now().UnixNano()
		a.call(content, err, request, writer)
		return
	}
	content.EndNano = time.Now().UnixNano()
	a.call(content, err, request, writer)
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
func (a *archiveWebsocketHandler) handler(msg []byte, writer http.ResponseWriter, request *http.Request) (content *ctx.Result, err error) {

	content = &ctx.Result{BeginNano: time.Now().UnixNano()}
	ms := strings.ToLower(request.URL.Path[a.arch.skip:])
	idx := strings.Index(ms, "/")
	if idx == -1 {
		content.Method = ms
		content.Version = a.arch.defaultMethodVersion[content.Method]
	} else {
		content.Method = ms[:idx]
		content.Version = a.arch.getVersion(content.Method, ms[idx+1:])
	}
	nc := newContext(writer, request)
	nc.Set("_message", string(msg))
	content.Request = a.createContext(nc)

	if err = runWare(content.Version, content.Request, a.arch.before); err != nil {
		content.EndNano = time.Now().UnixNano()
		return
	}
	re, err := a.arch.runMethod(content.Method, content.Version, content.Request)
	if err != nil {
		content.EndNano = time.Now().UnixNano()
		return
	}
	content.Response = re
	err = runWare(content.Version, content.Request, a.arch.after)
	return
}
