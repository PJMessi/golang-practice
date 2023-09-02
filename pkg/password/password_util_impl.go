package password

import (
	"regexp"
	"unicode"
)

type UtilImpl struct {
	Util
}

func NewUtil() Util {
	return &UtilImpl{}
}

func (utility *UtilImpl) IsStrong(plainPassword string) (bool, error) {
	// password should be at least 8 characters long and must have at least 1 lowercase character, 1 uppercase character,
	// 1 digit and 1 special character (!@#$%^&*()_+{})

	if len(plainPassword) < 8 {
		return false, nil
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range plainPassword {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\|-]`).MatchString(string(char)):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial, nil
}
