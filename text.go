package toolkit

import (
	"fmt"
)

var _randChars string

/*
func randChars() string {
	if len(_randChars) == 0 {
		alphabets := "abcdefghijklmnopqrstuvwxyz"
		alphabetsCap := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numerics := "0123456789"
		whitespaces := "!@#$^&*-_+"
		_randChars = alphabets + numerics + alphabetsCap + whitespaces
	}
	return _randChars
}

func SetRandChars(chars string) string {
	_randChars = chars
	return _randChars
}
*/

func RandomString(length int) string {
	return GenerateRandomString("", length)
}

func Sprintf(pattern string, parms ...interface{}) string {
	return fmt.Sprintf(pattern, parms...)
}

func Printf(pattern string, parms ...interface{}) {
	fmt.Printf(pattern, parms...)
}

func Println(s ...interface{}) {
	fmt.Println(s)
}
