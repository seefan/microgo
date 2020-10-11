package httpserver

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/seefan/microgo/ctx"
	"net"
	"net/http"
	"strings"
)

type archiveHandler struct {
	arch          *Archive
	call          func(interface{}, error, *http.Request, http.ResponseWriter)
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
		a.call(result, err, request, writer)
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

type archiveWebsocketHandler struct {
	arch          *Archive
	call          func(interface{}, error) []byte
	createContext func(httpContext *HTTPContext) ctx.Entry
}

func (a *archiveWebsocketHandler) write(conn net.Conn, op ws.OpCode, msg string) {

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
func (a *archiveWebsocketHandler) handler(msg []byte, writer http.ResponseWriter, request *http.Request) (result interface{}, err error) {
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
	nc.Set("_message", string(msg))
	c := a.createContext(nc)

	if err = runWare(version, c, a.arch.before); err != nil {
		return
	}
	if result, err = a.arch.runMethod(method, version, c); err != nil {
		return
	}
	err = runWare(version, c, a.arch.after)
	return
}
