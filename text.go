package toolkit

import (
	"fmt"
	"strings"
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

func Split(txt string, splitterChars []string) (splitValues []string,
	splitters []string) {
	txtLen := len(txt)
	tmp := ""
	for i := 0; i < txtLen; i++ {
		c := string(txt[i])
		if !HasMember(splitterChars, c) {
			tmp += c
		} else {
			splitValues = append(splitValues, tmp)
			splitters = append(splitters, c)
			tmp = ""
		}
	}
	return
}

func SubStr(str string, start int, length int) string {
    if start < 0 {
        start = 0
    }
	if length == 0 {
		length = len(str)
	}
	strArr := strings.Split(str, "")
	i      := 0
	str     = ""
	for i=start; i<start+length; i++ {
		if i >= len(strArr) {
			break
		}
		str += strArr[i]
	}
	return str
}

func StrToLower(str string) string {
	return strings.ToLower(str)
}

func StrToUpper(str string) string {
	return strings.ToUpper(str)
}

func UcFirst(str string) string {
	first := SubStr(str, 0, 1)
	last  := SubStr(str, 1, 0)
	return strings.ToUpper(first)+strings.ToLower(last)
}

func UcWords(str string) string {
	strArr := strings.Split(str, " ")
	for i:=0; i<len(strArr); i++ {
		strArr[i] = UcFirst(strArr[i])
	}
	return strings.Join(strArr, " ")
}
