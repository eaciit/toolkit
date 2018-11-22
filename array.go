package toolkit

import (
	"reflect"
	"strings"
	"time"
	//"errors"
)

func VariadicToSlice(objs ...interface{}) *[]interface{} {
	result := []interface{}{}
	for _, v := range objs {
		result = append(result, v)
	}
	return &result
}

func MapToSlice(objects map[string]interface{}) []interface{} {
	results := make([]interface{}, 0)
	for _, v := range objects {
		results = append(results, v)
	}
	return results
}

func HasMember(g interface{}, find interface{}) bool {
	found := false
	if IsSlice(g) == false {
		return false
	}

	count := SliceLen(g)
	for i := 0; i < count; i++ {
		v := SliceItem(g, i)
		eq := Compare(v, find, "$eq")
		//Println("L:", v, " F:", find, ", B:", eq, reflect.ValueOf(find).Type().String())
		if eq {
			return true
		}
	}
	return found
}

func MemberIndex(g interface{}, find interface{}) (found bool, in int) {
	found = false
	if IsSlice(g) == false {
		in = -1
		return
	}

	count := SliceLen(g)
	for in = 0; in < count; in++ {
		v := SliceItem(g, in)
		eq := Compare(v, find, "$eq")
		if eq {
			found = true
			return
		}
	}

	in = -1
	return
}

func ToInterfaceArray(o interface{}) []interface{} {
	if IsSlice(o) == false {
		return []interface{}{}
	}

	//Printf("Slice Data: %s\n", JsonString(o))
	var ret []interface{}
	for i := 0; i < SliceLen(o); i++ {
		sliceItem := SliceItem(o, i)
		//Printf("%d Item: %s\n", i, JsonString(sliceItem))
		ret = append(ret, sliceItem)
	}
	return ret
}

func Compare(v1 interface{}, v2 interface{}, op string) bool {
	if !strings.HasPrefix(op, "$") {
		op = "$" + op
	}

	vv1 := reflect.Indirect(reflect.ValueOf(v1))
	vv2 := reflect.Indirect(reflect.ValueOf(v2))

	t1 := strings.ToLower(TypeName(v1))
	t2 := strings.ToLower(TypeName(v2))

	if strings.Contains(t1, "int") || strings.Contains(t1, "float") || strings.Contains(t1, "number") {
		//--- is a number
		// lets convert all to float64 for simplicity
		var vv1o, vv2o float64
		vv1o = ToFloat64(v1, 10, RoundingAuto)
		vv2o = ToFloat64(v2, 10, RoundingAuto)

		if op == "$eq" {
			return vv1o == vv2o
		} else if op == "$ne" {
			return vv1o != vv2o
		} else if op == "$lt" {
			return vv1o < vv2o
		} else if op == "$lte" {
			return vv1o <= vv2o
		} else if op == "$gt" {
			return vv1o > vv2o
		} else if op == "$gte" {
			return vv1o >= vv2o
		}
	} else if strings.Contains(t1, "time.time") || strings.Contains(t2, "time.time") {
		//--- is a time.Time
		vv1o := time.Now()
		if !strings.Contains(t1, "time.time") {
			vv1o, _ = time.Parse(time.RFC3339, v1.(string))
		} else {
			vv1o = vv1.Interface().(time.Time)
		}
		vv2o := vv2.Interface().(time.Time)
		if op == "$eq" {
			return vv1o.Sub(vv2o).Nanoseconds() == 0
		} else if op == "$ne" {
			return vv1o != vv2o
		} else if op == "$lt" {
			return vv1o.Before(vv2o)
		} else if op == "$lte" {
			return vv1o == vv2o || vv1o.Before(vv2o)
		} else if op == "$gt" {
			return vv1o.After(vv2o)
		} else if op == "$gte" {
			return vv1o == vv2o || vv1o.After(vv2o)
		}

	} else if strings.Contains(t1, "bool") {
		vv1o := vv1.Bool()

		if vv2.Kind() == reflect.Bool {
			vv2o := vv2.Bool()
			if op == "$eq" {
				return vv1o == vv2o
			} else if op == "$ne" {
				return vv1o != vv2o
			}
		} else {
			return false
		}
	} else {
		//--- will be string
		vv1o := ToString(vv1.Interface())
		vv2o := ToString(vv2.Interface())
		if op == "$eq" {
			return vv1o == vv2o
		} else if op == "$ne" {
			return vv1o != vv2o
		} else if op == "$lt" {
			return strings.Compare(vv1o, vv2o) < 0
		} else if op == "$lte" {
			return strings.Compare(vv1o, vv2o) <= 0
		} else if op == "$gt" {
			return strings.Compare(vv1o, vv2o) > 0
		} else if op == "$gte" {
			return strings.Compare(vv1o, vv2o) >= 0
		}
	}

	return false
}
