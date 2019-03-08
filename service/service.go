package service

// Service base service
type Service interface {
	//service version
	Version() string
	//Based on this url
	Path() string
}
