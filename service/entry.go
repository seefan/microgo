package service

// Entry common param
type Entry interface {
	Get(string) string
	Value(string) Value
	GetSlice(string) []string
}
