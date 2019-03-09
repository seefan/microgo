package main

import (
	"errors"
	"github.com/seefan/microgo/ctx"
	"github.com/seefan/microgo/httpserver"
	"github.com/seefan/microgo/run"
	"github.com/seefan/microgo/service"
	"github.com/seefan/microgo/test"
)

func main() {
	s := httpserver.NewHTTPServer("localhost", 8889)
	s.Prefix = "/svr"
	s.Register(&test.TestService{}).Before(service.Ware{Next: func(entry ctx.Entry) (err error) {
		name := entry.String("name")
		if name != "jack" {
			err = errors.New("not login")
		}
		return
	}})
	//if err := s.Start(context.Background()); err != nil {
	//	println(err.Error())
	//}
	run.Run(s)

}
