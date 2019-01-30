package httpserver

import (
	"net/http"
	"net/url"
)

// HTTPContext context
type HTTPContext struct {
	forms    url.Values
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
		c.forms = request.Form
	}
	return c
}

// Get get on param
func (h *HTTPContext) Get(name string) string {
	if vs, ok := h.forms[name]; ok {
		if len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}

// GetSlice get slice
func (h *HTTPContext) GetSlice(name string) []string {
	if vs, ok := h.forms[name]; ok {
		return vs
	}
	return nil
}
