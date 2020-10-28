package template

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/seefan/goerr"
)

//Template 模板
type Template struct {
	tpl    *template.Template
	Cached bool
	Root   string
	ext    string
}

//New 创建模板
func New(root, ext string) (t *Template, err error) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	t = &Template{Cached: false, Root: root, ext: ext}
	t.init()
	err = t.makeDir(t.tpl, t.Root)
	return
}
func (t *Template) init() {
	t.tpl = template.New("")
	t.tpl.Delims("{{", "}}")
	t.tpl.Funcs(Func)
}

//MakeFile 生成一个文件
func (t *Template) MakeFile(src string, w io.Writer, param interface{}) error {
	if !t.Cached {
		t.init()
		if err := t.makeDir(t.tpl, t.Root); err != nil {
			return err
		}
	}
	//不存在的文件，不处理
	tm := t.tpl.Lookup(src)
	if tm == nil {
		return goerr.String("%s not found", src)
	}

	if err := tm.Execute(w, param); err != nil {
		return goerr.Errorf(err, "run %s error", src)
	}
	return nil
}
func (t *Template) makeDir(tpl *template.Template, path string) error {
	//查看目录下文件
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if strings.HasSuffix(path, "/") {
		path += "/"
	}
	for _, fi := range fis {
		tmpPath := path + "/" + fi.Name()
		if fi.IsDir() { //如果是目录，就也目录下所有文件都拷贝过去
			if err := t.makeDir(tpl, tmpPath); err != nil {
				return err
			}
		} else {
			extName := filepath.Ext(fi.Name())
			//特定扩展名的文件才可以作模板
			if t.ext == extName {
				if err = t.makeFile(tpl, tmpPath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (t *Template) makeFile(tpl *template.Template, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	if i := strings.Index(path, "/"); i != -1 {
		path = path[i+1:]
	}
	tp := tpl.New(path)
	_, err = tp.Parse(s)
	if err != nil {
		return err
	}
	return nil
}
