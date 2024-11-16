package libraryIO

import "strings"

func MakeSecret(str string) string {
	newString := make([]rune, len(str))
	for i := range str {
		newString[i] = '*'
	}
	return string(newString)
}

func makeSecret2(str string) string {
	newString := ""
	for range str {
		newString += "*"
	}
	return string(newString)
}

func makeSecret3(str string) string {
	return strings.Repeat("*", len(str))
}

func makeSecret4(str string) string {
	if len(str) == 1 {
		return "*"
	}
	return "*" + makeSecret4(str[1:])
}
