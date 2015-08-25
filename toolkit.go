package toolkit

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"reflect"
	//"strconv"
	"fmt"
	"strings"
	"time"
)

type M map[string]interface{}

func (m M) Set(k string, v interface{}) M {
	m[k] = v
	return m
}

func (m M) Get(k string, d ...interface{}) interface{} {
	if get, b := m[k]; b {
		return get
	} else {
		if len(d) > 0 {
			return d[0]
		} else {
			return nil
		}
	}
}

func (m M) GetInt(k string) int {
	i := m.Get(k, 0)
	switch i.(type) {
	case int:
		return i.(int)
	case int32, int64:
		return int(i.(int64))
	case float32:
		return int(i.(float32))
	case float64:
		return int(i.(float64))
	default:
		return 0
	}
}

func (m M) GetFloat32(k string) float32 {
	i := m.Get(k, 0)
	switch i.(type) {
	case int, int32, int64:
		return float32(i.(int))
	case float32:
		return i.(float32)
	case float64:
		return float32(i.(float64))
	default:
		return 0
	}
}

func (m M) GetFloat64(k string) float64 {
	i := m.Get(k, 0)
	switch i.(type) {
	case int, int32, int64:
		return float64(i.(int))
	case float32:
		return float64(i.(float32))
	case float64:
		return i.(float64)
	default:
		return 0
	}
}

func (m M) Has(k string) bool {
	_, has := m[k]
	return has
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

func MakeDate(layout string, value string) time.Time {
	t, e := time.Parse(layout, value)
	if e != nil {
		t, _ = time.Parse("2-Jan-2006", "1-Jan-1900")
		return t
	} else {
		return t
	}
}

func AddTime(dt0 time.Time, dt1 time.Time) time.Time {
	dtx := dt0
	return dtx.Add(dt1.Sub(MakeDate("03:04", "00:00")))
}

func Id(i interface{}) interface{} {
	//_ = "breakpoint"
	idFields := []interface{}{"_id", "ID", "Id", "id"}
	rv := reflect.ValueOf(i)

	//-- get key
	found := false
	var id interface{}
	if rv.Kind() == reflect.Map {
		mapkeys := rv.MapKeys()
		for _, mapkey := range mapkeys {
			idkey := mapkey.String()
			if HasMember(idFields, idkey) {
				idValue := rv.MapIndex(mapkey)
				if idValue.IsValid() {
					found = true
					id = idValue.Interface()
				}
			}
		}
	} else if rv.Kind() == reflect.Struct {
		for _, idkey := range idFields {
			idValue := rv.FieldByName(idkey.(string))
			if idValue.IsValid() {
				found = true
				id = idValue.Interface()
			}
		}
	} else if rv.Kind() == reflect.Ptr {
		elem := rv.Elem()
		for _, idkey := range idFields {
			idValue := elem.FieldByName(idkey.(string))
			if idValue.IsValid() {
				found = true
				id = idValue.Interface()
			}
		}
	} else {
		//_ = "breakpoint"
		fmt.Printf("Kind: %s \n", rv.Kind().String())
	}

	if found {
		return id
	} else {
		return nil
	}
}

func Value(i interface{}, fieldName string, def interface{}) interface{} {
	rv := reflect.ValueOf(i)
	var ret interface{}
	found := false
	if rv.Kind() == reflect.Map {
		mapkeys := rv.MapKeys()
		for i := 0; i < len(mapkeys) && !found; i++ {
			mapkey := mapkeys[i]
			mapkeyname := mapkey.String()
			if mapkeyname == fieldName {
				found = true
				mapvalue := rv.MapIndex(mapkey)
				if mapvalue.IsNil() {
					ret = def
				} else {
					ret = mapvalue.Interface()
				}
			}
		}
	} else if rv.Kind() == reflect.Struct {
		fv := rv.FieldByName(fieldName)
		if fv.IsValid() {
			found = true
			if (fv.Kind() == reflect.Struct || fv.Kind() == reflect.Map) && fv.IsNil() {
				ret = def
			} else {
				ret = fv.Interface()
			}
		}
	}

	if !found {
		return def
	} else {
		return ret
	}
}

func Field(o interface{}, fieldName string) (reflect.Value, bool) {
	ref := reflect.ValueOf(o)
	if !ref.IsValid() {
		return ref, false
	}
	es := ref.Elem()
	fi := es.FieldByName(fieldName)
	if fi.IsValid() {
		return fi, true
	}
	return fi, false
}

func JsonString(o interface{}) string {
	bs, e := json.Marshal(o)
	if e != nil {
		return "{}"
	}
	return string(bs)
}

func ObjFromString(s string, result interface{}) error {
	b := []byte(s)
	e := json.Unmarshal(b, result)
	return e
}

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

func PathDefault(removeSlash bool) string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//dir, _ := os.Getwd()
	if removeSlash == false {
		dir = dir + "/"
	}
	return dir
}

func DecodeByte(bytesData []byte, result interface{}) error {
	buf := bytes.NewBuffer(bytesData)
	dec := gob.NewDecoder(buf)
	e := dec.Decode(result)
	return e
}

func GetEncodeByte(obj interface{}) []byte {
	b, e := EncodeByte(obj)
	if e != nil {
		return new(bytes.Buffer).Bytes()
	}
	return b
}

func EncodeByte(obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	gw := gob.NewEncoder(buf)
	err := gw.Encode(obj)
	if err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

func GetIP() ([]string, error) {

	ret := make([]string, 0)
	he := func(err error) ([]string, error) {
		return ret, err
	}

	ifaces, err := net.Interfaces()
	he(err)

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		he(err)

		for _, addr := range addrs {
			interfaceTxt := addr.String()
			if strings.HasSuffix(interfaceTxt, "24") {
				interfaceTxt = interfaceTxt[0 : len(interfaceTxt)-3]
				ret = append(ret, interfaceTxt)
			}
		}
	}
	if len(ret) == 0 {
		ret = append(ret, "127.0.0.1")
	}
	return ret, nil
}
