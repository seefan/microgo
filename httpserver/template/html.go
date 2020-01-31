package template

type HTML struct {
	URL     string
	Context HTMLContext
}
type HTMLContext map[string]interface{}

func (c HTMLContext) Title(title string) {
	c["title"] = title
}
