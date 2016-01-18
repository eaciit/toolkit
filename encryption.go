package toolkit

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(baseChars string, n int) string {
	if baseChars == "" {
		baseChars = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz-_"
	}
	baseCharsLen := len(baseChars)

	rnd := ""
	for i := 0; i < n; i++ {
		x := RandInt(baseCharsLen)
		rnd += string(baseChars[x])
	}
	return rnd
}
