package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/seefan/microgo/service"
	"log"
	"net/http"
	"strconv"

	"github.com/seefan/microgo/server"
)

// HTTPServer for basic function
type HTTPServer struct {
	server.Server
	svr    *http.Server
	isRun  bool
	header map[string]string
	svc    map[string]*archive
}

// NewHTTPServer create new http server
func NewHTTPServer(host string, port int) *HTTPServer {
	hs := &HTTPServer{
		header: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "content-type",
		},
		svc: make(map[string]*archive),
	}
	hs.Server.Init(host, port)
	return hs
}

// Start the server
func (h *HTTPServer) Stop() error {
	return nil
}

// Start the server
func (h *HTTPServer) Start(ctx context.Context) (err error) {
	if h.InitFunc != nil {
		h.InitFunc()
	}
	return h.run(ctx)
}
func (h *HTTPServer) Register(svc service.Service) {
	if a, ok := h.svc[svc.Name()]; ok {
		a.Put(svc)
	} else {
		h.svc[svc.Name()] = NewArchive()
		h.svc[svc.Name()].Put(svc)
	}
}
func (h *HTTPServer) run(ctx context.Context) (err error) {
	h.svr = &http.Server{Addr: h.Address()}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var result []interface{}
		var err error
		defer func() {
			re := make(map[string]interface{})
			if err != nil {
				re["error"] = err.Error()
			}
			if result != nil {
				for i, v := range result {
					if v != nil {
						if e, ok := v.(error); ok {
							re["error"] = e.Error()
						} else if i == 0 {
							re["data"] = v
						} else {
							re["data"+strconv.Itoa(i)] = v
						}
					}
				}
				if _, ok := re["error"]; !ok {
					re["error"] = 0
				}
				if bs, err := json.Marshal(re); err == nil {
					if _, err := writer.Write(bs); err != nil {
						log.Println(err)
					}
				}
			}
		}()
		meta, err := GetMetaFromURL(request.URL.Path)
		if err != nil {
			return
		}
		sv, ok := h.svc[meta.Service]
		if !ok {
			err = errors.New("UnknownService")
			return
		}
		svc := sv.Get(meta.Version)

		for k, v := range h.header {
			writer.Header().Add(k, v)
		}

		result, err = svc.RunMethod(meta.Method, newContext(writer, request))
	})
	h.isRun = true
	log.Println("http server is start")
	err = h.svr.ListenAndServe()
	if err != nil {
		h.isRun = false
	}
	ctx.Done()
	return
}
