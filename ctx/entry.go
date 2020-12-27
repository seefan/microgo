package ctx

// Entry common param
type Entry interface {
	String(string) string
	Value(string) Value
	Body() string
}

type Result struct {
	Response  interface{}
	Request   Entry
	BeginNano int64
	EndNano   int64
	Method    string
	Version   string
}
