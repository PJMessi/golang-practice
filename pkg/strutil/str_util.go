package strutil

import (
	"bytes"
	"unicode"
)

func ConvertToBytes(str string) []byte {
	return []byte(str)
}

func IsEmpty(str string) bool {
	return str == ""
}

func ConvertFromBytes(bytes []byte) string {
	return string(bytes)
}

func PascalCaseToCamelCase(s string) string {
	var result bytes.Buffer
	var prevIsUpper, currIsUpper bool

	for i, r := range s {
		currIsUpper = unicode.IsUpper(r)

		if i == 0 && currIsUpper {
			// Make the first character lowercase.
			result.WriteRune(unicode.ToLower(r))
		} else if i > 0 && !prevIsUpper && currIsUpper {
			// Insert an underscore before the uppercase character.
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}

		prevIsUpper = currIsUpper
	}

	return result.String()
}
