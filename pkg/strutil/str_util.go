package strutil

func ConvertToBytes(str string) []byte {
	return []byte(str)
}

func IsEmpty(str string) bool {
	return str == ""
}

func ConvertFromBytes(bytes []byte) string {
	return string(bytes)
}
