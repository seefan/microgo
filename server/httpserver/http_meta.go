package httpserver

import (
	"errors"
	"strings"
)

// HTTPMeta service
type HTTPMeta struct {
	Service string
	Version string
	Method  string
}

// GetMetaFromURL get meta from url
func GetMetaFromURL(url string) (*HTTPMeta, error) {
	us := strings.Split(url, "/")
	if len(us) == 4 {
		return &HTTPMeta{us[1], us[2], us[3]}, nil
	}
	return nil, errors.New("WrongURL")
}
