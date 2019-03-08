package httpserver

import "errors"

var Skip = 0

// HTTPMeta service
type HTTPMeta struct {
	Service string
	Version string
	Method  string
}

// GetMetaFromURL get meta from url
func GetMetaFromURL(url string) (*HTTPMeta, error) {
	pos := make([]int, 2)
	idx := 0
	size := len(url)
	if url[size-1] == '/' {
		size--
	}
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
	return &HTTPMeta{url[:pos[1]], url[pos[1]+1 : pos[0]], url[pos[0]+1 : size]}, nil
}
