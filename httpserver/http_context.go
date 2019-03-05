package httpserver

import (
	"github.com/seefan/microgo/service"
	"net/http"
)

// HTTPContext context
type HTTPContext struct {
	service.Context
	Request  *http.Request
	Response http.ResponseWriter
}

// NewContext new NewContext
func newContext(writer http.ResponseWriter, request *http.Request) *HTTPContext {
	c := &HTTPContext{
		Request:  request,
		Response: writer,
	}
	if request.ParseForm() == nil {
		c.Context = *service.NewContext(request.Form)
	}
	return c
}
