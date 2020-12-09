package httpserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/seefan/goerr"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/httpserver/template"
	"github.com/seefan/microgo/server"
	"github.com/seefan/microgo/service"
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
	//web socket prefix
	websocketURL string
	//common header
	header map[string]string
	//path:service
	arch map[string]*Archive

	//method context
	Context func(*HTTPContext) ctx.Entry
	//Marshal format
	Marshal func(result *ctx.Result, err error) ([]byte, error)
	//result format
	Result func(result *ctx.Result, err error) interface{}
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
	hs.Result = func(result *ctx.Result, err error) interface{} {
		re := make(map[string]interface{})
		if err != nil {
			hs.RuntimeLog(err)
			if hs.Debug {
				re["error"] = goerr.Error(err).Trace()
			} else {
				re["error"] = err.Error()
			}

		} else if result.Data != nil {
			if e, ok := result.Data.(error); ok && e != nil {
				re["error"] = e.Error()
			} else {
				re["data"] = result.Data
				re["error"] = 0
			}
		} else {
			re["error"] = 0
		}
		return re
	}
	hs.Marshal = func(result *ctx.Result, err error) ([]byte, error) {
		re := hs.Result(result, err)
		return json.Marshal(re)
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
func (h *HTTPServer) Start() (err error) {
	if h.InitFunc != nil {
		h.InitFunc()
	}
	return h.run()
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
	sg := service.NewService(svc)
	path := h.getPath(sg.Path())
	if _, ok := h.arch[path]; !ok {
		h.arch[path] = newArchive(path)
	}

	h.arch[path].put(sg)
	return h.arch[path]
}

//RegisterAfterWare only register ware
func (h *HTTPServer) RegisterAfterWare(svc service.Service, md ...service.Ware) {
	sg := service.NewService(svc)
	path := h.getPath(sg.Path())
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

//SetWebsocketURL websocket url prefix
func (h *HTTPServer) SetWebsocketURL(url string) {
	h.websocketURL = url
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

//SetTemplatePath static files
func (h *HTTPServer) SetTemplatePath(path, ext string, cached ...bool) error {
	h.templatePath = path
	tpl, err := template.New(path, ext)
	if err != nil {
		return err
	}
	if len(cached) > 0 && cached[0] == true {
		tpl.Cached = true
	}
	h.tpl = tpl
	return nil
}

//RegisterBeforeWare only register ware
func (h *HTTPServer) RegisterBeforeWare(svc service.Service, md ...service.Ware) {
	sg := service.NewService(svc)
	path := h.getPath(sg.Path())
	if _, ok := h.arch[path]; ok {
		for _, m := range md {
			h.arch[path].Before(m)
		}
	}
}
func (h *HTTPServer) html(ht *template.HTML, err error, request *http.Request, w io.Writer) {
	if err != nil {
		h.RuntimeLog(err)
		return
	}
	if h.tpl != nil {
		if len(request.Form) > 0 {
			rspForm := make(map[string]interface{})
			for k, v := range request.Form {
				if len(v) == 0 {
					rspForm[k] = ""
				} else if len(v) == 1 {
					rspForm[k] = v[0]
				} else {
					rspForm[k] = v
				}
			}
			ht.Context["_form"] = rspForm
		}
		if err := h.tpl.MakeFile(ht.URL, w, ht.Context); err != nil {
			h.RuntimeLog(err)
		}
	}
}
func (h *HTTPServer) run() error {
	mux := http.NewServeMux()

	h.svr = &http.Server{Addr: h.Address(), Handler: mux}
	if h.docOnline != "" {
		mux.HandleFunc(h.docOnline, h.handleDoc)
	}

	for path, s := range h.arch {
		if h.websocketURL != "" && strings.HasPrefix(path, h.websocketURL) {
			mux.Handle(path, &archiveWebsocketHandler{arch: s, createContext: h.Context, call: func(result *ctx.Result, err error) []byte {
				bs, err := h.Marshal(result, err)
				if err != nil {
					h.RuntimeLog(err)
					return nil
				}
				return bs
			}})
		} else {
			mux.Handle(path, &archiveHandler{arch: s, createContext: h.Context, call: func(result *ctx.Result, err error, request *http.Request, writer http.ResponseWriter) {
				for k, v := range h.header {
					writer.Header().Set(k, v)
				}
				if r, ok := result.Data.(*template.HTML); ok {
					writer.Header().Set("Content-Type", "text/html;charset=utf-8")
					h.html(r, err, request, writer)
				} else {
					bs, err := h.Marshal(result, err)
					if err != nil {
						h.RuntimeLog(err)
						return
					}
					if _, err := writer.Write(bs); err != nil {
						h.RuntimeLog(err)
					}
				}
			}})
		}
	}

	if h.staticPath != "" {
		mux.Handle(h.staticURL, http.StripPrefix(h.staticURL, http.FileServer(http.Dir(h.staticPath))))
	}

	h.isRun = true
	err := h.svr.ListenAndServe()
	if err != nil {
		h.isRun = false
	}
	return err
}
