package toolkit

import (
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
)

var gobs []string

func RegisterGobObject(o interface{}) {
	name := reflect.ValueOf(o).Type().Name()
	if HasMember(gobs, name) {
		return
	}

	gob.Register(o)
	gobs = append(gobs, name)
}

func TypeName(o interface{}) string {
	v := reflect.ValueOf(o)
	return v.Type().String()
	/*
		v := reflect.ValueOf(o)
		if !v.IsValid() {
			return ""
		}

		var t reflect.Type
		if v.Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
			if v.Kind() == reflect.Ptr {
				return "*" + TypeName(v.Elem().Interface())
			} else {
				t = v.Type()
			}
		} else {
			t = reflect.TypeOf(o)
		}

		//t := v.Type()
		name := t.Name()
		pkg := t.PkgPath()
		if pkg != "" {
			return pkg + "." + name
		} else {
			return name
		}
	*/
}

func IsNilOrEmpty(x interface{}) bool {
	rv := reflect.Indirect(reflect.ValueOf(x))
	return !rv.IsValid() || x == reflect.Zero(reflect.TypeOf(x)).Interface()
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
		return nil, errors.New("GetEmptySliceElement: " + TypeName(o) + " is not a slice")
	}
	sliceType := rv.Type().Elem()
	newelem := reflect.New(sliceType)
	//fmt.Println(sliceType.String())
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

func MapKeys(o interface{}) []interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Map {
		return []interface{}{}
	}

	var ret []interface{}
	keyvalues := v.MapKeys()
	for _, kv := range keyvalues {
		ret = append(ret, kv.Interface())
	}
	return ret
}

func MapLen(o interface{}) int {
	return len(MapKeys(o))
}

func MapItem(o interface{}, key interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Map {
		return nil
	}
	item := v.MapIndex(reflect.ValueOf(key))
	return item.Interface()
}

func SliceItem(o interface{}, index int) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		return nil
	}
	if v.Len()-1 < index {
		return nil
	}
	return v.Index(index).Interface()
}

func Serde(o interface{}, dest interface{}, serdeType string) error {
	bs, e := ToBytesWithError(o, serdeType)
	if len(bs) == 0 {
		return errors.New("toolkit.Serde: Serialization Fail " + e.Error())
	}

	e = FromBytes(bs, serdeType, dest)
	if e != nil {
		return errors.New("toolkit.Serde: Deserialization fail " + e.Error())
	}

	return nil
}
