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
