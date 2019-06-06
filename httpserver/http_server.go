package httpserver

import (
	"context"
	"encoding/json"
	"github.com/seefan/goerr"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/server"
	"github.com/seefan/microgo/service"
	"io"
	"net/http"
	"strings"
)

// HTTPServer for basic function
type HTTPServer struct {
	server.Server
	isRun bool
	//Debug output debug info
	Debug bool
	//doc path
	docOnline string
	//static file
	staticPath string
	//static file
	staticDir string
	//common prefix
	Prefix string
	//server
	svr *http.Server
	//mux
	//mux *http.ServeMux

	//common header
	header map[string]string
	//path:service
	arch map[string]*archive

	//method context
	Context func(*HTTPContext) ctx.Entry
	//Marshal format
	Marshal func(result interface{}, err error, w io.Writer)
	//result format
	Result func(result interface{}, err error) interface{}
	//log
	RuntimeLog func(err error)
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
		RuntimeLog: func(err error) {
			//do nothing
		},
	}
	hs.Result = func(result interface{}, err error) interface{} {
		re := make(map[string]interface{})
		if err != nil {
			hs.RuntimeLog(err)
			if hs.Debug {
				re["error"] = goerr.Error(err).Trace()
			} else {
				re["error"] = err.Error()
			}

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
	}
	hs.Marshal = func(result interface{}, err error, w io.Writer) {
		re := hs.Result(result, err)
		if bs, err := json.Marshal(re); err == nil {
			if _, err := w.Write(bs); err != nil {
				hs.RuntimeLog(err)
			}
		} else {
			hs.RuntimeLog(err)
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
	if h.isRun {
		h.isRun = false
		h.svr.Shutdown()
		if h.CloseFunc != nil {
			h.CloseFunc()
		}
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
func (h *HTTPServer) getPath(p string) (path string) {
	if strings.HasPrefix(p, "/") {
		path = h.Prefix + p
	} else {
		path = h.Prefix + "/" + p
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return
}
func (h *HTTPServer) Register(svc service.Service) *archive {
	path := h.getPath(svc.Path())
	if _, ok := h.arch[path]; !ok {
		h.arch[path] = NewArchive()
	}
	h.arch[path].Put(svc)
	return h.arch[path]
}

// only register ware
func (h *HTTPServer) RegisterAfterWare(svc service.Service, md ...service.Ware) {
	path := h.getPath(svc.Path())
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

func (h *HTTPServer) SetStaticPath(path, dir string) {
	h.staticPath = path
	h.staticDir = dir
}

// only register ware
func (h *HTTPServer) RegisterBeforeWare(svc service.Service, md ...service.Ware) {
	path := h.getPath(svc.Path())
	if _, ok := h.arch[path]; ok {
		for _, m := range md {
			h.arch[path].Before(m)
		}
	}
}

//ServeHTTP server http method
func (h *HTTPServer) serve(writer http.ResponseWriter, request *http.Request, sv *archive, path string) {
	if strings.ToLower(request.Method) == "options" {
		writer.WriteHeader(204)
		return
	}
	for k, v := range h.header {
		writer.Header().Set(k, v)
	}
	var result interface{}
	var err error

	//meta, err := getMetaFromURL(request.URL.Path)
	//if err != nil {
	//	return
	//}
	defer func() {
		h.Marshal(result, err, writer)
	}()
	var method, version string
	ms := request.URL.Path[len(path):]
	idx := strings.Index(ms, "/")
	if idx == -1 {
		method = ms
	} else {
		method = ms[:idx]
		version = ms[idx+1:]
	}

	nc := newContext(writer, request)
	c := h.Context(nc)
	if err = runWare(version, c, sv.before); err != nil {
		return
	}
	if result, err = sv.runMethod(method, version, c); err != nil {
		return
	}
	if err = runWare(version, c, sv.after); err != nil {
		return
	}
}
func (h *HTTPServer) run(ctx context.Context) error {

	mux := http.NewServeMux()
	h.svr = &http.Server{Addr: h.Address(), Handler: mux}
	if h.docOnline != "" {
		mux.HandleFunc(h.docOnline, h.handleDoc)
	}
	if h.staticDir != "" {
		mux.Handle(h.staticPath, http.FileServer(http.Dir(h.staticDir)))
	}
	//http.HandleFunc(h.Prefix+"/", h.ServeHTTP)
	for path, s := range h.arch {
		println("register", path)
		mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			h.serve(writer, request, s, path)
		})
	}
	h.isRun = true
	err := h.svr.ListenAndServe()
	if err != nil {
		h.isRun = false
	}
	ctx.Done()
	return err
}
