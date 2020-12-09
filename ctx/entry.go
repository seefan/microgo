package ctx

// Entry common param
type Entry interface {
	String(string) string
	Value(string) Value
}

type Result struct {
	Data      interface{}
	BeginNano int64
	EndNano   int64
}
