package service

// Entry common param
type Entry interface {
	Get(string) string
	GetSlice(string) []string
}
