package httpserver

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/seefan/microgo/ctx"
)

//HTTPContext context
type HTTPContext struct {
	form     url.Values
	Request  *http.Request
	Response http.ResponseWriter
}

//NewContext new NewContext
func newContext(writer http.ResponseWriter, request *http.Request) *HTTPContext {
	c := &HTTPContext{
		Request:  request,
		Response: writer,
	}
	contentType := request.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if request.ParseForm() == nil {
			c.form = request.Form
		}
	}
	return c
}

//SetForm Set form param
func (h *HTTPContext) SetForm(forms url.Values) {
	for k, v := range forms {
		h.form[k] = v
	}
}

//String  get string param
func (h *HTTPContext) String(name string) string {
	if vs, ok := h.form[name]; ok {
		if len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}

//Value get param value
func (h *HTTPContext) Value(name string) ctx.Value {
	if vs, ok := h.form[name]; ok {
		if len(vs) > 0 {
			return ctx.Value(vs[0])
		}
	}
	return ""
}

//Set set param
func (h *HTTPContext) Set(name, value string) {
	h.form[name] = []string{value}
}

//Slice get param slice
func (h *HTTPContext) Slice(name string) []string {
	return h.form[name]
}

//Keys list key
func (h *HTTPContext) Keys() []string {
	var ks []string
	for k := range h.form {
		ks = append(ks, k)
	}
	return ks
}
