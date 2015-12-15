package toolkit

import (
	"strings"
)

func ToBytes(data interface{}, encoderId string) []byte {
	encoderId = strings.ToLower(encoderId)
	if encoderId == "" {
		encoderId = "json"
	}

	if encoderId == "json" {
		return Jsonify(data)
	} else if encoderId == "gob" {
		b, e := EncodeByte(data)
		if e != nil {
			return []byte{}
		} else {
			return b
		}
	}
	return []byte{}
}
