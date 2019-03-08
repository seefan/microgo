package test

import (
	"github.com/seefan/microgo/ctx"
)

type TestService struct {
}

func (TestService) Hello(entry ctx.Entry) interface{} {
	name := entry.String("name")
	return "hello " + name
}
func (TestService) Path() string {
	return "test"
}
func (TestService) Version() string {
	return "1.0"
}
