package toolkit

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type M map[string]interface{}
type Ms []M

var ErrorPathNotFound = errors.New("Path requested is not available")

func (m M) Set(k string, v interface{}) M {
	m[k] = v
	return m
}

func (m M) Sets(vs ...interface{}) M {
	k := ""
	for i, v := range vs {
		if i%2 == 0 {
			k, _ = v.(string)
		} else {
			if k != "" {
				m.Set(k, v)
				k = ""
			}
		}
	}
	return m
}

func (m M) PathSet(k string, v interface{}, pathDelim string) M {
	doms := strings.Split(k, pathDelim)
	//Logger().Debugf("key:%s doms:%s", k, JsonString(doms))
	if len(doms) == 1 {
		//Logger().Debugf("Set %s to %v", k, v)
		m.Set(doms[0], v)
	} else if len(doms) > 1 {
		map0 := M{}
		dom0 := m.Get(doms[0])
		if reflect.Indirect(reflect.ValueOf(dom0)).Kind() == reflect.Map {
			map0, _ = ToM(dom0)
		}
		map0.PathSet(strings.Join(doms[1:], pathDelim), v, pathDelim)
		//Logger().Debugf("Iterate child %s to v", k, v)
		m.Set(doms[0], map0)
	}
	//Logger().Debugf("M: %s", JsonString(m))
	return m
}

func (m M) Get(k string, d ...interface{}) interface{} {
	if get, b := m[k]; b {
		if IsNilOrEmpty(get) && len(d) > 0 {
			get = d[0]
		}
		return get
	} else {
		if len(d) > 0 {
			return d[0]
		} else {
			return nil
		}
	}
}

func (m M) GetRef(k string, d, out interface{}) {
	/*
		defer func() {
			if r := recover(); r != nil {
			}
		}()
	*/
	valget := reflect.Indirect(reflect.ValueOf(m.Get(k, d)))
	valout := reflect.ValueOf(out)
	valout.Elem().Set(valget)
}

const (
	CaseAsIs  string = ""
	CaseUpper        = "upper"
	CaseLower        = "lower"
)

var (
	DefaultCase = CaseAsIs
)

func ToMCase(data interface{}, casePattern string) (M, error) {
	return tom(data, casePattern)
}

func ToM(data interface{}) (M, error) {
	return tom(data, DefaultCase)
}

func ToMTag(data interface{}, tagName string) (M, error) {
	return tomTagName(data, DefaultCase, tagName)
}

var _tagName string

func TagName() string {
	if _tagName == "" {
		_tagName = "json"
	}
	return _tagName
}

func SetTagName(name string) {
	_tagName = name
}

func tom(data interface{}, namePattern string) (M, error) {
	return tomTagName(data, namePattern, TagName())
}

func tomTagName(data interface{}, namePattern string, tagName string) (M, error) {
	rv := reflect.Indirect(reflect.ValueOf(data))
	// Create emapty map as a result
	res := M{}

	// Because of the difference behaviour of Struct type and Map type, we need to check the data element type
	if rv.Kind() == reflect.Struct {
		// Iterate through all the available field
		for i := 0; i < rv.NumField(); i++ {
			// Get the field type
			f := rv.Type().Field(i)
			fieldName := f.Name
			fieldTagNames := strings.Split(f.Tag.Get(tagName), ",")
			/*
				if len(fieldTagNames) > 1 && fieldTagNames[1] == "omitempty" &&  f {
					continue
				}
			*/
			fieldTagName := fieldTagNames[0]
			if fieldTagName != "" {
				fieldName = fieldTagName
			}
			if fieldName == "-" {
				continue
			}

			switch namePattern {
			case CaseLower:
				fieldName = strings.ToLower(fieldName)

			case CaseUpper:
				fieldName = strings.ToUpper(fieldName)
			}

			// If the type is struct but not time.Time or is a map
			kind := f.Type.Kind()
			if kind == reflect.Ptr {
				kind = f.Type.Elem().Kind()
			} else if (kind == reflect.Struct && f.Type != reflect.TypeOf(time.Time{})) || kind == reflect.Map {
				// Then we need to call this function again to fetch the sub value
				subRes, err := tomTagName(rv.Field(i).Interface(), namePattern, tagName)
				if err != nil {
					return nil, err
				}
				res[fieldName] = subRes
				// Skip the rest
				continue
			} else if kind == reflect.Slice {
				slice := rv.Field(i)
				count := slice.Len()
				elemType := slice.Type().Elem().Kind()
				switch elemType {
				case reflect.Map, reflect.Struct, reflect.Ptr, reflect.Slice:
					ms := make([]M, count)
					for si := 0; si < count; si++ {
						data := slice.Index(si)
						datam, err := tomTagName(data.Interface(), namePattern, tagName)
						if err != nil {
							return nil, err
						}
						ms[si] = datam
					}
					res[fieldName] = ms

				default:
					res[fieldName] = slice.Interface()
				}
				continue
			}

			// If the type is time.Time or is not struct and map then put it in the result directly
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						//-- currently do nothing
					}
				}()
				res[fieldName] = rv.Field(i).Interface()
			}()
		}

		// Return the result
		return res, nil
	} else if rv.Kind() == reflect.Map {
		// If the data element is kind of map
		// Iterate through all avilable keys
		for _, key := range rv.MapKeys() {

			// Get the map value type of the specified key
			if elem := reflect.Indirect(rv.MapIndex(key)); elem.IsValid() {
				t := elem.Type()
				// If the type is struct but not time.Time or is a map
				if (t.Kind() == reflect.Struct && t != reflect.TypeOf(time.Time{})) || t.Kind() == reflect.Map {
					// Then we need to call this function again to fetch the sub value
					subRes, err := tom(rv.MapIndex(key).Interface(), namePattern)
					if err != nil {
						return nil, err
					}
					res[key.String()] = subRes

					// Skip the rest
					continue
				}
			}
			// If the type is time.Time or is not struct and map then put it in the result directly
			res[key.String()] = rv.MapIndex(key).Interface()
		}

		// Return the result
		return res, nil
	}

	// If the data element is not map or struct then return error
	return nil, Errorf("Expecting struct or map object but got %s", rv.Kind())
}

func (m M) ToBytes(encodertype string, others ...interface{}) []byte {
	encodertype = strings.ToLower(encodertype)
	if encodertype == "" {
		encodertype = "json"
	}
	/*
		if encodertype == "json" {
			bs, e := json.Marshal(m)
			if e != nil {
				return []byte{}
			}
			return bs
		}
	*/
	return ToBytes(m, encodertype)
	//return []byte{}
}

func (m *M) Cast(k string, d interface{}) error {
	var e error
	if m.Has(k) == false {
		return fmt.Errorf("No data for key %s", k)
	}
	b, e := json.Marshal(m.Get(k, nil))
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, d)
	return e
}

func (m M) GetBytes(k string) []byte {
	bs, err := base64.StdEncoding.DecodeString(m.GetString(k))
	if err != nil {
		return []byte{}
	}
	return bs
}

func (m M) GetFloat64(k string) float64 {
	i := m.Get(k, 0)
	return ToFloat64(i, 6, RoundingAuto)
}

func (m M) GetString(k string) string {
	s := m.Get(k, "")
	return ToString(s)
}

func (m M) GetInt(k string) int {
	i := m.Get(k, 0)
	return ToInt(i, RoundingAuto)
}

func (m *M) Unset(k string) {
	delete(*m, k)
}

func (m M) GetFloat32(k string) float32 {
	i := m.Get(k, 0)
	return ToFloat32(i, 6, RoundingAuto)
}

func (m M) GetBool(k string) bool {
	b := strings.ToLower(m.GetString(k))
	return b == "1" || b == "true" || b == "y" || b == "t" || b == "yes"
}

func (m M) Has(k string) bool {
	_, has := m[k]
	return has
}

func (m M) PathGet(path string) (interface{}, error) {
	pathlist := strings.Split(path, ".")
	var curobj interface{} = m
	for _, nextpath := range pathlist {
		switch curobj.(type) {
		case M:
			curM := curobj.(M)
			var found bool
			curobj, found = curM[nextpath]
			if !found {
				return nil, ErrorPathNotFound
			}

		case map[string]interface{}:
			curM := curobj.(map[string]interface{})
			var found bool
			curobj, found = curM[nextpath]
			if !found {
				return nil, ErrorPathNotFound
			}

		default:
			return nil, ErrorPathNotFound
		}
	}

	return curobj, nil
}

func (m M) Keys() []string {
	var ret []string
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
}

func (m M) Values() []interface{} {
	var ret []interface{}
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func CopyM(from, to *M,
	copyFieldIfNotExist bool,
	exceptFields []string) {
	//Printf("Copy from: %s to %s\n", JsonString(from), JsonString(to))
	var exceptFieldsIface []interface{}
	Serde(exceptFields, &exceptFieldsIface, "")
	fromm := *from
	tom := *to
	for f, fv := range fromm {
		if !HasMember(exceptFieldsIface, f) {
			if tom.Has(f) {
				tom.Set(f, fv)
			} else if copyFieldIfNotExist {
				tom.Set(f, fv)
			}
		}
	}
	*to = tom
}

func (m M) Merge(from M, overwrite bool) {
	for k, v := range from {
		_, ok := m[k]
		if ok && !overwrite {
			continue
		}
		m[k] = v
	}
}
