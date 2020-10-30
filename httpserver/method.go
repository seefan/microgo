package httpserver

import "github.com/seefan/microgo/ctx"

//Method method map
type Method = interface{}

// //MethodInputPutput method
type MethodInputPutput = func(entry ctx.Entry) interface{}

// //MethodOnlyInput no output
type MethodOnlyInput = func(entry ctx.Entry)

// //MethodOnlyOutput no input
type MethodOnlyOutput = func() interface{}

// //MethodNotInputOutput no input no output
type MethodNoInputOutput = func()
