package toolkit

import (
	"reflect"
)

func IsValid(o interface{}) bool {
	//t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if !v.IsValid() {
		return false
	}
	if v.IsNil() {
		return false
	}
	return true
}

func IfEq(has interface{}, want interface{}, a interface{}, b interface{}) interface{} {
	if has == want {
		return a
	} else {
		return b
	}
}

func IfNe(has interface{}, dontWant interface{}, a interface{}, b interface{}) interface{} {
	if has == dontWant {
		return a
	} else {
		return b
	}
}

func IfFn(f func() bool, a interface{}, b interface{}) interface{} {
	if f() {
		return a
	} else {
		return b
	}
}
