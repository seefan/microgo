package template

import (
	"html/template"
	"strings"
	"time"
)

//TemplateFuncs template funcs
var Func template.FuncMap = make(map[string]interface{})

const (
	fmtDateTimeString = "2006-01-02 15:04:05"
	fmtDateString     = "2006-01-02"
)

func init() {
	//format date
	Func["fmt_date"] = func(s time.Time) string {
		return s.Format(fmtDateString)
	}
	Func["fmt_datetime"] = func(s time.Time) string {
		return s.Format(fmtDateTimeString)
	}
	Func["start_with"] = func(s string, d string) bool {
		return strings.HasPrefix(s, d)
	}
	Func["html"] = rawHTML
	Func["attr"] = rawHTMLAttr
	Func["js"] = rawJS
	Func["css"] = rawCSS
	Func["url"] = rawURL

	// UtilFuncs["dict_extends"] = func(d model.Dict, i string) string {
	// 	switch i {
	// 	case "6":
	// 		return d.Param6
	// 	case "1":
	// 		return d.Param1
	// 	case "2":
	// 		return d.Param2
	// 	case "3":
	// 		return d.Param3
	// 	case "4":
	// 		return d.Param4
	// 	case "5":
	// 		return d.Param5
	// 	default:
	// 		return d.Param0
	// 	}
	// }
	//TemplateFuncs["equals"] = equals
	//TemplateFuncs["notequals"] = func(ss ...interface{}) bool {
	//	return !equals(ss...)
	//}
	Func["imgsrc"] = func(img, defaultImg string) string {
		if len(img) > 0 {
			return img
		} else {
			return defaultImg
		}
	}

}

//func equals(ss ...interface{}) bool {
//	for i, s := range ss {
//		if i == 0 {
//			continue
//		}
//		if utils.AsString(s) != utils.AsString(ss[i-1]) {
//			return false
//		}
//	}
//	return true
//}

func rawHTML(text string) template.HTML {
	return template.HTML(text)
}
func rawHTMLAttr(text string) template.HTMLAttr {
	return template.HTMLAttr(text)
}

func rawCSS(text string) template.CSS {
	return template.CSS(text)
}

func rawJS(text string) template.JS {
	return template.JS(text)
}

func rawJSStr(text string) template.JSStr {
	return template.JSStr(text)
}

func rawURL(text string) template.URL {
	return template.URL(text)
}
