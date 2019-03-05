package service

// Entry common param
type Entry interface {
	Get(string) Value
	Value(string) Value
	GetSlice(string) []Value
}
