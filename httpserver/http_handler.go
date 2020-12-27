package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/seefan/goerr"
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
	archMap map[string]*Archive
	path    func(string) (string, string)
	message func(string) (map[string]string, error)
	call    func(*ctx.Result, error) []byte

	createContext func(httpContext *HTTPContext) ctx.Entry
}

//ServeHTTP server http method
func (a *archiveWebsocketHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(request, writer)
	if err != nil {
		content := &ctx.Result{BeginNano: time.Now().UnixNano()}
		bs := a.call(content, err)
		wsutil.WriteServerMessage(conn, ws.OpText, bs)
		return
	}

	go func() {
		defer conn.Close()
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			content := &ctx.Result{BeginNano: time.Now().UnixNano()}
			if err != nil {
				bs := a.call(content, err)
				if err := wsutil.WriteServerMessage(conn, ws.OpText, bs); err != nil {
					a.call(content, err)
					break
				}
				continue
			}
			if op == ws.OpPing {
				_ = wsutil.WriteServerMessage(conn, ws.OpPong, []byte{})
				continue
			}
			err = a.handler(content, msg, writer, request)
			if err != nil {
				bs := a.call(content, err)
				if err := wsutil.WriteServerMessage(conn, ws.OpText, bs); err != nil {
					a.call(content, err)
					break
				}
				continue
			}
			bs := a.call(content, nil)
			err = wsutil.WriteServerMessage(conn, op, bs)
			if err != nil {
				a.call(content, err)
				break
			}
		}
	}()

}
func (a *archiveWebsocketHandler) handler(content *ctx.Result, data []byte, writer http.ResponseWriter, request *http.Request) (err error) {
	path, msg := a.path(string(data))
	if len(path) < 2 {
		err = goerr.Errorf(goerr.String("Path:%s", path), "PathEmpty")
		content.EndNano = time.Now().UnixNano()
		return
	}
	ms := strings.ToLower(path)
	idx := strings.Index(ms[1:], "/")
	svr := string(ms[:idx+2])
	arch, exists := a.archMap[svr]
	if !exists {
		err = goerr.Errorf(goerr.String("Path:%s", path), "PathNotFound")
		content.EndNano = time.Now().UnixNano()
		return
	}
	param, err := a.message(msg)
	if err != nil {
		err = goerr.Errorf(err, "MessageError")
		content.EndNano = time.Now().UnixNano()
		return
	}
	ms = string(ms[idx+2:])
	idx = strings.Index(ms, "/")
	if idx == -1 {
		content.Method = ms
		content.Version = arch.defaultMethodVersion[content.Method]
	} else {
		content.Method = ms[:idx]
		content.Version = arch.getVersion(content.Method, ms[idx+1:])
	}
	nc := newContext(writer, request)
	nc.body = msg
	for k, v := range param {
		nc.Set(k, v)
	}
	content.Request = a.createContext(nc)

	if err = runWare(content.Version, content.Request, arch.before); err != nil {
		content.EndNano = time.Now().UnixNano()
		return
	}
	re, err := arch.runMethod(content.Method, content.Version, content.Request)
	if err != nil {
		content.EndNano = time.Now().UnixNano()
		return
	}
	content.Response = re
	err = runWare(content.Version, content.Request, arch.after)
	content.EndNano = time.Now().UnixNano()
	return
}
