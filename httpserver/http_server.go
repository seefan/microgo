package httpserver

import (
	"context"
	"encoding/json"
	"github.com/seefan/goerr"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/server"
	"github.com/seefan/microgo/service"
	"io"
	"log"
	"net/http"
	"strings"
)

// HTTPServer for basic function
type HTTPServer struct {
	server.Server
	svr   *http.Server
	isRun bool
	//common header
	header map[string]string
	//path:service
	arch map[string]*archive
	//doc path
	docOnline string
	//common prefix
	Prefix      string
	Context     func(*HTTPContext) ctx.Entry
	Output      func(result interface{}, err error, w io.Writer)
	BuildResult func(result interface{}, err error) interface{}
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
		docOnline: "/doc",
		arch:      make(map[string]*archive),
		Context: func(httpContext *HTTPContext) ctx.Entry {
			return httpContext
		},
		BuildResult: func(result interface{}, err error) interface{} {
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
			return re
		},
	}
	hs.Output = func(result interface{}, err error, w io.Writer) {
		re := hs.BuildResult(result, err)
		if bs, err := json.Marshal(re); err == nil {
			if _, err := w.Write(bs); err != nil {
				log.Println(err)
			}
		}
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
func (h *HTTPServer) Register(svc service.Service) *archive {
	var path string
	if strings.HasPrefix(svc.Path(), "/") {
		path = h.Prefix + svc.Path()
	} else {
		path = h.Prefix + "/" + svc.Path()
	}

	if _, ok := h.arch[path]; !ok {
		h.arch[path] = NewArchive()
	}
	h.arch[path].Put(svc)
	return h.arch[path]
}

// only register ware
func (h *HTTPServer) RegisterAfterWare(svc service.Service, md ...service.Ware) {
	path := h.Prefix + "/" + svc.Path()
	if _, ok := h.arch[path]; ok {
		for _, m := range md {
			h.arch[path].After(m)
		}
	}
}

// only register ware
func (h *HTTPServer) SetDocPath(path string) {
	h.docOnline = path
}

// only register ware
func (h *HTTPServer) RegisterBeforeWare(svc service.Service, md ...service.Ware) {
	path := h.Prefix + "/" + svc.Path()
	if _, ok := h.arch[path]; ok {
		for _, m := range md {
			h.arch[path].Before(m)
		}
	}
}
func (h *HTTPServer) run(ctx context.Context) error {
	h.svr = &http.Server{Addr: h.Address()}
	if h.docOnline != "" {
		http.HandleFunc(h.docOnline, h.handleDoc)
	}
	http.HandleFunc(h.Prefix+"/", func(writer http.ResponseWriter, request *http.Request) {
		if strings.ToLower(request.Method) == "options" {
			writer.WriteHeader(204)
			return
		}
		for k, v := range h.header {
			writer.Header().Set(k, v)
		}
		var result interface{}
		var err error

		meta, err := GetMetaFromURL(request.URL.Path)
		if err != nil {
			return
		}
		defer func() {
			putMeta(meta)
			h.Output(result, err, writer)
		}()
		sv, ok := h.arch[meta.Service]
		if !ok {
			err = goerr.Errorf(goerr.String("Service:%s Method:%s Version:%s", meta.Service, meta.Method, meta.Version), "UnknownService")
			return
		}

		nc := newContext(writer, request)
		c := h.Context(nc)
		if err = runWare(meta.Version, c, sv.before); err != nil {
			return
		}
		if result, err = sv.runMethod(meta.Method, meta.Version, c); err != nil {
			return
		}
		if err = runWare(meta.Version, c, sv.after); err != nil {
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
