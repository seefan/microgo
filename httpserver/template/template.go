package template

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/golangteam/function/errors"
	"github.com/golangteam/function/file"
)

//Template 模板
type Template struct {
	tpl    *template.Template
	Cached bool
	Root   string
}

//New 创建模板
func New(root string) *Template {
	t := &Template{Cached: false, Root: root}
	t.init()
	return t
}
func (t *Template) init() {
	t.tpl = template.New("")
	t.tpl.Delims("{{", "}}")
}

//MakeFile 生成一个文件
func (t *Template) MakeFile(src string, w io.Writer, param interface{}) error {
	src = t.Root + "/" + src
	//不存在的文件，不处理
	if file.FileIsNotExist(src) {
		return errors.New("%s not found", src)
	}
	//解析文件链接为真实路径
	input, err := filepath.EvalSymlinks(src)
	if err != nil {
		return errors.NewError(err, "%s is not found", src)
	}

	buf, err := ioutil.ReadFile(input)
	if err != nil {
		return errors.NewError(err, "%s is not readable", src)
	}
	tm := t.tpl.New(src)

	// add our funcmaps
	tm.Funcs(Func)
	// Bomb out if parse fails. We don't want any silent server starts.
	if _, err := tm.Parse(string(buf)); err != nil {
		return errors.NewError(err, "template %s parse error")
	}
	if err := t.tpl.ExecuteTemplate(w, input, param); err != nil {
		return errors.NewError(err, "run %s error", input)
	}
	return nil
}
