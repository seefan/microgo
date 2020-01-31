package httpserver

import (
	"context"
	"encoding/json"
	"github.com/seefan/goerr"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/httpserver/template"
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
	tpl   *template.Template
	isRun bool
	//Debug output debug info
	Debug bool
	//doc path
	docOnline string
	//static file
	staticPath string
	//static url
	staticURL string
	//template path
	templatePath string
	//common prefix
	Prefix string
	//server
	svr *http.Server
	//mux
	//mux *http.ServeMux

	//common header
	header map[string]string
	//path:service
	arch map[string]*Archive

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
		arch:      make(map[string]*Archive),
		Context: func(httpContext *HTTPContext) ctx.Entry {
			return httpContext
		},
		RuntimeLog: func(err error) {
			log.Println(err)
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

//Header Set common header
func (h *HTTPServer) Header(name, value string) {
	h.header[name] = value
}

//Stop stop the server
func (h *HTTPServer) Stop() error {
	if h.isRun {
		h.isRun = false
		if h.CloseFunc != nil {
			h.CloseFunc()
		}
		return h.svr.Close()
	}
	return nil
}

//Start the server
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

//Register register service
func (h *HTTPServer) Register(svc service.Service) *Archive {
	path := h.getPath(svc.Path())
	if _, ok := h.arch[path]; !ok {
		h.arch[path] = newArchive(path)
	}

	h.arch[path].Put(svc)
	return h.arch[path]
}

//RegisterAfterWare only register ware
func (h *HTTPServer) RegisterAfterWare(svc service.Service, md ...service.Ware) {
	path := h.getPath(svc.Path())
	if _, ok := h.arch[path]; ok {
		for _, m := range md {
			h.arch[path].After(m)
		}
	}
}

//SetDocPath only register ware
func (h *HTTPServer) SetDocPath(path string) {
	h.docOnline = path
}

//SetStaticPath static files
func (h *HTTPServer) SetStaticPath(path, url string) {
	h.staticPath = path
	if strings.HasSuffix(url, "/") || url == "" {
		h.staticURL = url
	} else {
		h.staticURL = url + "/"
	}
}

//SetStaticPath static files
func (h *HTTPServer) SetTemplatePath(path string) {
	h.templatePath = path
	h.tpl = template.New(path)

}

//RegisterBeforeWare only register ware
func (h *HTTPServer) RegisterBeforeWare(svc service.Service, md ...service.Ware) {
	path := h.getPath(svc.Path())
	if _, ok := h.arch[path]; ok {
		for _, m := range md {
			h.arch[path].Before(m)
		}
	}
}
func (h *HTTPServer) html(ht *template.HTML, err error, w io.Writer) {
	if err != nil {
		h.RuntimeLog(err)
		return
	}
	if h.tpl != nil {
		if err := h.tpl.MakeFile(ht.URL, w, ht.Context); err != nil {
			h.RuntimeLog(err)
		}
	}
}
func (h *HTTPServer) run(ctx context.Context) error {
	mux := http.NewServeMux()

	h.svr = &http.Server{Addr: h.Address(), Handler: mux}
	if h.docOnline != "" {
		mux.HandleFunc(h.docOnline, h.handleDoc)
	}

	for path, s := range h.arch {
		mux.Handle(path, &archiveHandler{arch: s, createContext: h.Context, call: func(result interface{}, err error, writer http.ResponseWriter) {
			for k, v := range h.header {
				writer.Header().Set(k, v)
			}
			if r, ok := result.(*template.HTML); ok {
				h.html(r, err, writer)
			} else {
				h.Marshal(result, err, writer)
			}
		}})
	}

	if h.staticPath != "" {
		mux.Handle(h.staticURL, http.StripPrefix(h.staticURL, http.FileServer(http.Dir(h.staticPath))))
	}

	h.isRun = true
	err := h.svr.ListenAndServe()
	if err != nil {
		h.isRun = false
	}
	ctx.Done()
	return err
}
