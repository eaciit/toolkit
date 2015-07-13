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
	"strings"
)

type M map[string]interface{}

func (m M) Set(k string, v interface{}) M {
	m[k] = v
	return m
}

func (m M) Get(k string, d interface{}) interface{} {
	if m.Has(k) == false {
		m[k] = d
	}
	return m[k]
}

func (m M) Has(k string) bool {
	_, has := m[k]
	return has
}

/**** NtbR *****/
/*
func SetField(o interface{}, fieldName string, value interface{}) {
	es := reflect.ValueOf(o).Elem()
	fi := es.FieldByName(fieldName)
	if fi.IsValid() {
		if fi.CanSet() {
			switch value.(type) {
			case int:
				x := value.(int)
				fi.SetInt(strconv.FormatInt(x, 10))

			case string:
				x := value.(string)
				fi.SetString(x)
			}
		}
	}
}
*/
func GetJsonString(o interface{}) string {
	bs, e := json.MarshalIndent(o, "", "\t")
	if e != nil {
		return "{}"
	}
	return string(bs)
}

func GetObjFromString(s string, result interface{}) error {
	b := []byte(s)
	e := json.Unmarshal(b, result)
	return e
}

func MapToSlice(objects map[string]interface{}) []interface{} {
	results := make([]interface{}, 0)
	for _, v := range objects {
		results = append(results, v)
	}
	return results
}

func GetField(o interface{}, fieldName string) (reflect.Value, bool) {
	es := reflect.ValueOf(o).Elem()
	fi := es.FieldByName(fieldName)
	if fi.IsValid() {
		return fi, true
	}
	return fi, false
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
