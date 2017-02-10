package toolkit

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
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

func FileChecksum(fileLocation string) string {
	f, err := os.Open(fileLocation)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return ""
	}

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return ""
	}

	hashed := hash.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}
