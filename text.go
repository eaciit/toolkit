package toolkit

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
