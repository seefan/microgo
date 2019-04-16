package httpserver

import "github.com/seefan/microgo/ctx"

type Method map[string]func(entry ctx.Entry) interface{}
