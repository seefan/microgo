package service

import "github.com/seefan/microgo/ctx"

type Ware struct {
	Next func(entry ctx.Entry) error
}
