package template

type HTML struct {
	URL     string
	Context map[string]interface{}
}
type HTMLParam struct {
	Form map[string]interface{}
	Context map[string]interface{}
}
//Html write html
func Html(url string, kvs ...interface{}) *HTML {
	html := &HTML{URL: url, Context: make(map[string]interface{})}
	for i := 0; i < len(kvs)-1; i += 2 {
		if k, ok := kvs[i].(string); ok {
			html.Context[k] = kvs[i+1]
		}
	}
	return html
}
