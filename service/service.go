package service

// Service base service
type Service interface {
	Version() string
	Name() string
}
