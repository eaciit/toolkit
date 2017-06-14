package toolkit

import (
	"encoding/gob"
	"errors"
	"reflect"
	"strings"
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
	//if x == nil {
	//	return true
	//}
	rv := reflect.Indirect(reflect.ValueOf(x))
	k := rv.Kind()
	if k == reflect.Slice {
		return false
	} else if k == reflect.String {
		if ToString(x) == "" {
			return true
		} else {
			return false
		}
	} else if k == reflect.Struct {
		return false
	} else if k == reflect.Bool {
		return false
	} else if strings.HasPrefix(k.String(), "int") || strings.Contains(k.String(), "float") {
		iszero := x == reflect.Zero(reflect.TypeOf(x)).Interface()
		if iszero {
			return true
		} else {
			return false
		}
	}

	invalid := !rv.IsValid()
	if invalid {
		return true
	}

	return rv.IsNil()
}

func IsNumber(o interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(o))
	ts := strings.ToLower(v.Type().String())
	if (strings.Contains(ts, "int") || strings.Contains(ts, "float")) && !strings.HasPrefix(ts, "interface{}") {
		return true
	}
	return false
}

func IsStringNumber(txt string, decsep string) (f float64, e error) {
	hasDes := false
	newtxt := "0"
	decPoint := 0
	//Printf("%v ", txt)
	for _, c := range txt {
		s := string(c)
		if strings.Compare(s, "0") >= 0 && strings.Compare(s, "9") <= 0 {
			newtxt += s
			if hasDes {
				decPoint += 1
			}
		} else if s == decsep {
			if !hasDes {
				newtxt += "."
				hasDes = true
			} else {
				e = errors.New("IsStringNumber: Multiple decimal separator found")
				return
			}
		} else {
			//Printfn("%v %v", txt, s)
			e = errors.New("IsStringNumber: Wrong character " + txt)
			return
		}
	}
	if strings.HasSuffix(newtxt, ".") {
		newtxt += "0"
	}
	//Printfn("%v",newtxt)
	f = ToFloat64(newtxt, decPoint, RoundingAuto)
	return
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
	//Println(newelem.Type().String())
	if string(sliceType.String()[0]) == "*" {
		return Value2Interface(newelem), nil
	} else {
		return Value2Interface(newelem.Elem()), nil
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
	//fmt.Printf("Type: %s \n", t.String())
	return reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Interface()
}

func SliceLen(o interface{}) int {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		return 0
	}
	return v.Len()
}

func SliceSubset(o interface{}, lbound, ubound int) interface{} {
	v := reflect.Indirect(reflect.ValueOf(o))
	l := v.Len()
	if lbound < l && ubound < l {
		var arrays reflect.Value
		for i := lbound; i <= ubound; i++ {
			elem := v.Index(i)
			if i == lbound {
				arrays = reflect.MakeSlice(elem.Type(), 0, 0)
			}
			arrays = reflect.Append(arrays, elem)
		}
		return arrays.Interface()
	}
	return nil
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
	return Value2Interface(v.Index(index))
}

func SliceSetItem(o interface{}, i int, d interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(o))
	if v.Kind() != reflect.Slice {
		err := "SliceSetItem: object returned is not a slice. It is " + v.Kind().String() + " " + v.Type().String()
		//Println(err)
		return errors.New(err)
	}
	currentLen := v.Len()
	if i >= currentLen {
		//i = currentLen + 1
		//Println("Set capacity to ", i+1)
		//v.SetCap(i + 1)
		v1 := reflect.Append(v, reflect.ValueOf(d))
		v.Set(v1)
	} else {
		v.Index(i).Set(reflect.ValueOf(d))
	}
	return nil
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

func Value2Interface(vi reflect.Value) interface{} {
	vik := vi.Type().String()
	if strings.Contains(vik, "string") {
		return vi.String()
	} else if (vik == "int" || vik == "int8" || vik == "int16" || vik == "int32" || strings.Contains(vik, "uint")) && !strings.Contains(vik, "interface") {
		return int(vi.Int())
	} else if strings.Contains(vik, "float") {
		return vi.Float()
	} else if strings.Contains(vik, "bool") {
		return vi.Bool()
	} else {
		return vi.Interface()
	}
}

func ExecFunc(fn interface{}, ins ...interface{}) (outs []reflect.Value, e error) {
	rvfn := reflect.ValueOf(fn)
	if rvfn.Kind() != reflect.Func {
		e = errors.New("Execfunc: First parameter should be a function")
		return
	}
	var rvins []reflect.Value
	for _, in := range ins {
		rvins = append(rvins, reflect.ValueOf(in))
	}
	outs = rvfn.Call(rvins)
	return
}
