package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/seefan/microgo/server"
	"github.com/seefan/microgo/service"
	"log"
	"net/http"
	"strings"
)

// HTTPServer for basic function
type HTTPServer struct {
	server.Server
	svr    *http.Server
	isRun  bool
	header map[string]string
	arch   map[string]*archive
	prefix string
}

// NewHTTPServer create new http server
func NewHTTPServer(host string, port int) *HTTPServer {
	hs := &HTTPServer{
		header: map[string]string{
			"Content-Type":                 "application/json;charset=UTF-8",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
			"Access-Control-Allow-Headers": "Cache-Control, Pragma, Origin, Authorization, Content-Type, X-Requested-With",
		},
		arch: make(map[string]*archive),
	}
	hs.Server.Init(host, port)
	return hs
}

// Set common header
func (h *HTTPServer) Header(name, value string) {
	h.header[name] = value
}

// Start the server
func (h *HTTPServer) Stop() error {
	if h.CloseFunc != nil {
		h.CloseFunc()
	}
	if h.isRun {
		return h.svr.Close()
	}
	return nil
}

// Start the server
func (h *HTTPServer) Start(ctx context.Context) (err error) {
	if h.InitFunc != nil {
		h.InitFunc()
	}
	return h.run(ctx)
}
func (h *HTTPServer) Register(svc service.Service, md ...service.Ware) {
	if _, ok := h.arch[svc.Name()]; !ok {
		h.arch[svc.Name()] = NewArchive()
	}
	h.arch[svc.Name()].Put(svc)
	h.arch[svc.Name()].BeginWare(svc.Name(), md...)
}
func (h *HTTPServer) RegisterEndWard(svc string, md ...service.Ware) {
	if _, ok := h.arch[svc]; ok {
		h.arch[svc].EndWare(svc, md...)
	}
}
func (h *HTTPServer) RegisterBeginWard(svc string, md ...service.Ware) {
	if _, ok := h.arch[svc]; ok {
		h.arch[svc].BeginWare(svc, md...)
	}
}
func (h *HTTPServer) run(ctx context.Context) error {
	h.svr = &http.Server{Addr: h.Address()}

	http.HandleFunc(h.prefix+"/", func(writer http.ResponseWriter, request *http.Request) {
		if strings.ToLower(request.Method) == "options" {
			writer.WriteHeader(204)
			return
		}
		for k, v := range h.header {
			writer.Header().Set(k, v)
		}
		var result interface{}
		var err error
		defer func() {
			re := make(map[string]interface{})
			if err != nil {
				re["error"] = err.Error()
			} else if result != nil {
				if e, ok := result.(error); ok && e != nil {
					re["error"] = e.Error()
				} else {
					re["data"] = result
					re["error"] = 0
				}
			} else {
				re["error"] = 0
			}
			if bs, err := json.Marshal(re); err == nil {
				if _, err := writer.Write(bs); err != nil {
					log.Println(err)
				}
			}
		}()
		meta, err := GetMetaFromURL(request.URL.Path)
		if err != nil {
			return
		}
		sv, ok := h.arch[meta.Service]
		if !ok {
			err = errors.New("UnknownService")
			return
		}
		svc := sv.Get(meta.Version)
		c := newContext(writer, request)
		if err = sv.RunWare(svc.Name, c, sv.begin); err != nil {
			return
		}
		if result, err = svc.RunMethod(meta.Method, c); err != nil {
			return
		}
		if err = sv.RunWare(svc.Name, c, sv.end); err != nil {
			return
		}
	})
	h.isRun = true
	err := h.svr.ListenAndServe()
	if err != nil {
		h.isRun = false
	}
	ctx.Done()
	return err
}
