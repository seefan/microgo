package httpserver

import "testing"

func TestGetMetaFromURL(t *testing.T) {
	urls := []string{"/service/hello/say/1.0", "/hello/say/", "//service/hello/say/1.0", "/service/hello/say/"}
	for _, url := range urls {
		if m, e := GetMetaFromURL(url); e != nil {
			t.Error(e)
		} else {
			t.Log("service=", m.Service, "version=", m.Version, "method=", m.Method)
		}
	}
}
