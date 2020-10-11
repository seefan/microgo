package template

import (
	"fmt"
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
	Func["fmt_date_string"] = func(date string) string {
		return date[:10]
	}
	Func["fmt_datetime"] = func(s time.Time) string {
		return s.Format(fmtDateTimeString)
	}
	Func["start_with"] = func(s string, d string) bool {
		return strings.HasPrefix(s, d)
	}
	Func["left"] = func(s string, size int) string {
		if len(s)<=size{
			return s
		}
		return s[:size]+"..."
	}
	Func["rawHTML"] = rawHTML
	Func["hidden"] = func(id,value string) template.HTML {
		return template.HTML(fmt.Sprintf(`<input type="hidden" id="%s" name="%s" data-bind="%s" value="%s"/>`, id, id, id, value))
	}
	Func["mv"] = func(id string,m map[string]interface{}) interface{} {
		if v,ok:=m[id];ok{
			if s,ok:=v.(string);ok{
				return s
			}
			return fmt.Sprint(m[id])
		}
		return ""
	}
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
	Func["default"] = func(str,defStr string) string {
		if str!=""{
			return str
		}
		return defStr
	}
	Func["add"] = func(n1,n2 int) int {
		return n1+n2
	}
	Func["dec"] = func(n1,n2 int) int {
		return n1-n2
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
