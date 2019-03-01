package test

import (
	"github.com/seefan/microgo/service"
)

type TestService struct {
}

func (TestService) Hello(entry service.Entry) interface{} {
	name := entry.Get("name")
	return "hello " + name
}
func (TestService) Name() string {
	return "test"
}
func (TestService) Version() string {
	return "1.0"
}
