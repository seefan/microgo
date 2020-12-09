package test

import (
	"math/rand"
	"strconv"
	"time"

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
	time.Sleep(time.Millisecond)
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
	ts := rand.Intn(1000)
	time.Sleep(time.Millisecond * time.Duration(ts))
	println(ts)
	name := entry.String("name")
	return "hello " + name
}

func (TestService1) Version() string {
	return "1.1"
}
