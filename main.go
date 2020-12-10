package main

import (
	"errors"

	"github.com/golangteam/function/run"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/httpserver"
	"github.com/seefan/microgo/service"
	"github.com/seefan/microgo/test"
)

func main() {

	run.Run(func() run.Runnable {
		s := httpserver.NewHTTPServer("localhost", 8889)
		println("httpserver start ", "localhost", 8889)
		//s.SetTemplatePath("/Volumes/doc/test/tpl", ".html")
		//s.Prefix = "/svr"
		s.Result = func(result *ctx.Result, err error) interface{} {
			re := make(map[string]interface{})
			if err != nil {
				re["error"] = err.Error()
			} else if result.Response != nil {
				if e, ok := result.Response.(error); ok && e != nil {
					re["error"] = e.Error()
				} else {
					re["data"] = result.Response
					re["error"] = 0
				}
			} else {
				re["error"] = 0
			}
			re["ts"] = (result.EndNano - result.BeginNano) / 1000000
			return re
		}
		s.Register(&test.TestService{}).Before(service.Ware{Next: func(entry ctx.Entry) (err error) {
			name := entry.String("name")
			if name != "jack" {
				err = errors.New("not login")
			}
			return
		}})
		s.Register(&test.TestService1{})
		return s
	})
}
