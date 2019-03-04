package main

import (
	"github.com/seefan/microgo/httpserver"
	"github.com/seefan/microgo/run"
	"github.com/seefan/microgo/test"
)

func main() {
	s := httpserver.NewHTTPServer("localhost", 8889)
	s.Register(&test.TestService{})
	//if err := s.Start(context.Background()); err != nil {
	//	println(err.Error())
	//}
	run.Run(s)

}
