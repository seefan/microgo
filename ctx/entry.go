package ctx

// Entry common param
type Entry interface {
	String(string) string
	Value(string) Value
}
