package toolkit

import (
	"encoding/json"
	"fmt"
	"strings"
)

type M map[string]interface{}
type Ms []M

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

func ToM(v interface{}) (M, error) {
	bs, e := json.Marshal(v)
	if e != nil {
		return M{}, fmt.Errorf("Unable to cast to M : " + e.Error())
	}

	m := M{}
	e = json.Unmarshal(bs, &m)
	if e != nil {
		return m, fmt.Errorf("Unable to uncast to M from bytes: " + e.Error())
	}

	return m, nil
}

func (m M) ToBytes(encodertype string, encodeFn interface{}) []byte {
	encodertype = strings.ToLower(encodertype)
	if encodertype == "" {
		encodertype = "json"
	}
	if encodertype == "json" {
		bs, e := json.Marshal(m)
		if e != nil {
			return []byte{}
		}
		return bs
	}
	return []byte{}
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

func (m M) GetString(k string) string {
	s := m.Get(k, "")
	return s.(string)
}

func (m M) GetInt(k string) int {
	i := m.Get(k, 0)
	return ToInt(i)
}

func (m *M) Unset(k string) {
	delete(*m, k)
}

func (m M) GetFloat32(k string) float32 {
	i := m.Get(k, 0)
	return ToFloat32(i)
}

func (m M) Has(k string) bool {
	_, has := m[k]
	return has
}
