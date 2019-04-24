package httpserver

import (
	"errors"
	"sync"
)

var (
	Skip     = 0
	poolMeta = &sync.Pool{
		New: func() interface{} {
			return new(HTTPMeta)
		},
	}
)

// HTTPMeta service
type HTTPMeta struct {
	Service string
	Version string
	Method  string
}

func putMeta(meta *HTTPMeta) {
	poolMeta.Put(meta)
}

// GetMetaFromURL get meta from url
func getMetaFromURL(url string) (*HTTPMeta, error) {
	pos := make([]int, 2)
	idx := 0
	size := len(url)
	for i := size - 1; i >= 0; i-- {
		if url[i] == '/' {
			pos[idx] = i
			idx++
		}

		if idx == 2 {
			break
		}
	}
	if idx != 2 {
		return nil, errors.New("WrongURL")
	}
	m := poolMeta.Get().(*HTTPMeta)
	m.Service = url[:pos[1]]
	m.Version = url[pos[0]+1 : size]
	m.Method = url[pos[1]+1 : pos[0]]
	return m, nil
}
