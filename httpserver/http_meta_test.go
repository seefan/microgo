package httpserver

import "testing"

func TestGetMetaFromURL(t *testing.T) {
	urls := []string{"/service/hello/1.0/say", "/hello/1.0/say", "//service/hello/1.0/say", "/service/hello/1.0/say/"}
	for _, url := range urls {
		if m, e := GetMetaFromURL(url); e != nil {
			t.Error(e)
		} else {
			t.Log("service=", m.Service, "version=", m.Version, "method=", m.Method)
		}

	}
}
