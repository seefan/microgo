package main

import (
	"github.com/seefan/microgo/httpserver"
	"github.com/seefan/microgo/run"
	"github.com/seefan/microgo/server"
	"github.com/seefan/microgo/test"
)

func main() {

	run.Run(func() server.Runnable {
		s := httpserver.NewHTTPServer("localhost", 8889)
		println("httpserver start ", "localhost", 8889)
		//s.SetTemplatePath("/Volumes/doc/test/tpl", ".html")
		//s.Prefix = "/svr"
		// s.Register(&test.TestService{}).Before(service.Ware{Next: func(entry ctx.Entry) (err error) {
		// 	// name := entry.String("name")
		// 	// if name != "jack" {
		// 	// 	err = errors.New("not login")
		// 	// }
		// 	return
		// }})
		s.Register(&test.TestService1{})
		return s
	})
}
