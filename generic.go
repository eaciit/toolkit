package toolkit

import (
	"reflect"
)

func TypeName(o interface{}) string {
	v := reflect.ValueOf(o)
	if !v.IsValid() {
		return ""
	}
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	t := v.Type()
	name := t.Name()
	pkg := t.PkgPath()
	if pkg != "" {
		return pkg + "." + name
	} else {
		return name
	}
}
