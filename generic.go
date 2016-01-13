package toolkit

import (
	"errors"
	"fmt"
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

func IsPointer(o interface{}) bool {
	v := reflect.ValueOf(o)
	return v.Kind() == reflect.Ptr
}

func IsSlice(o interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(o))
	return v.Kind() == reflect.Slice
}

func GetEmptySliceElement(o interface{}) (interface{}, error) {
	rv := reflect.Indirect(reflect.ValueOf(o))
	if rv.Kind() != reflect.Slice {
		return nil, errors.New("GetEmptySliceElement: " + rv.Kind().String() + " is not a slice")
	}
	sliceType := rv.Type().Elem()
	newelem := reflect.New(sliceType)
	if string(sliceType.String()[0]) == "*" {
		return newelem.Interface(), nil
	} else {
		return newelem.Elem().Interface(), nil
	}
}

func AppendSlice(o interface{}, v interface{}) error {
	rvPtr := reflect.ValueOf(o)
	if rvPtr.Kind() != reflect.Ptr {
		return errors.New("AppendSlice: Is not a pointer of slice")
	}
	rv := reflect.Indirect(rvPtr)
	if rv.Kind() != reflect.Slice {
		return errors.New("AppendSlice: " + rv.Kind().String() + " is not a pointer of slice")
	}
	rv = reflect.Append(rv, reflect.ValueOf(v))
	rvPtr.Elem().Set(rv)
	return nil
}

func MakeSlice(o interface{}) interface{} {
	t := reflect.TypeOf(o)
	fmt.Printf("Type: %s \n", t.String())
	return reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Interface()
}

func SliceLen(o interface{}) int {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		return 0
	}
	return v.Len()
}

func SliceItem(o interface{}, index int) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		return nil
	}
	if v.Len()-1 < index {
		return nil
	}
	return v.Index(index)
}
