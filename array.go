package toolkit

import (
//"reflect
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

func HasMember(g []interface{}, find interface{}) bool {
	found := false
	for _, v := range g {
		if v == find {
			return true
		}
	}
	return found
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
