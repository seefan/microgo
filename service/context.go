package service

type Context struct {
	c map[string][]string
}

// NewContext new NewContext
func NewContext() *Context {
	return &Context{
		c: make(map[string][]string),
	}
}
func (h *Context) Set(name, value string) {
	h.c[name] = []string{value}
}
func (h *Context) SetSlice(name string, values []string) {
	h.c[name] = values
}

// Get get on param
func (h *Context) Get(name string) string {
	if vs, ok := h.c[name]; ok {
		if len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}
func (h *Context) Value(name string) Value {
	if vs, ok := h.c[name]; ok {
		if len(vs) > 0 {
			return Value(vs[0])
		}
	}
	return ""
}

// GetSlice get slice
func (h *Context) GetSlice(name string) []string {
	if vs, ok := h.c[name]; ok {
		return vs
	}
	return nil
}
