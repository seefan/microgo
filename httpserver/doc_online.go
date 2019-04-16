package httpserver

import (
	"net/http"
)

var (
	html string
)

func (h *HTTPServer) handleDoc(writer http.ResponseWriter, request *http.Request) {
	for path, ms := range h.arch {
		for name, method := range ms.method {
			writer.Write([]byte("<br>"))
			writer.Write([]byte(path))
			writer.Write([]byte("/"))
			writer.Write([]byte(name))
			writer.Write([]byte("/(version) ["))
			for version := range method {
				writer.Write([]byte(" "))
				writer.Write([]byte(version))
				writer.Write([]byte(" "))
			}
			writer.Write([]byte("]"))
		}
	}
}
