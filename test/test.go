package test

import (
	"github.com/seefan/microgo/ctx"
	"strconv"
)

type TestService struct {
}

func (TestService) HelloWorld(entry ctx.Entry) interface{} {

	return "hello "
}
func (TestService) Hello(entry ctx.Entry) interface{} {
	name := entry.String("name")
	a := 3
	b := 3 / (a - 3)
	return "hello " + name + strconv.Itoa(b)
}
func (TestService) Path() string {
	return "test"
}
func (TestService) Version() string {
	return "1.0"
}

type TestService1 struct {
}

func (TestService1) Hello(entry ctx.Entry) interface{} {
	name := entry.String("name")
	a := 3
	b := 3 / (a - 3)
	return "hello " + name + strconv.Itoa(b)
}
func (TestService1) Path() string {
	return "test"
}
func (TestService1) Version() string {
	return "1.1"
}
