package main

import (
	"context"
	"github.com/seefan/microgo/server/httpserver"
	"github.com/seefan/microgo/test"
)

func main() {
	s := httpserver.NewHTTPServer("localhost", 8888)
	s.Register(&test.TestService{})
	if err := s.Start(context.Background()); err != nil {
		println(err.Error())
	}

}
