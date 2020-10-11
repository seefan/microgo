package test

import (
	"strconv"

	"github.com/seefan/microgo/ctx"
)

type TestService struct {
}

// func (TestService) Default(entry ctx.Entry) interface{} {
// 	c := make(template.HTMLContext)
// 	c.Title("test it")
// 	return &template.HTML{URL: "important.txt", Context: c}
// }
func (TestService) HelloWorld(entry ctx.Entry) interface{} {
	return "hello "
}
func (TestService) Hello(entry ctx.Entry) interface{} {
	name := entry.String("name")
	a := 3
	b := 3 / (a - 4)
	return "hello " + name + strconv.Itoa(b)
}
func (TestService) Path() string {
	return ""
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

func (TestService1) Version() string {
	return "1.1"
}
