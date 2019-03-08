package httpserver

import (
	"github.com/seefan/microgo/ctx"
	"net/http"
	"net/url"
)

// HTTPContext context
type HTTPContext struct {
	form     url.Values
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
		c.form = request.Form
	}
	return c
}

// Set form param
func (h *HTTPContext) SetForm(forms url.Values) {
	for k, v := range forms {
		h.form[k] = v
	}
}

//  get string param
func (h *HTTPContext) String(name string) string {
	if vs, ok := h.form[name]; ok {
		if len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}

// get param value
func (h *HTTPContext) Value(name string) ctx.Value {
	if vs, ok := h.form[name]; ok {
		if len(vs) > 0 {
			return ctx.Value(vs[0])
		}
	}
	return ""
}

// get param slice
func (h *HTTPContext) Slice(name string) []string {
	return h.form[name]
}
