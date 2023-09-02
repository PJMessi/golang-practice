package password

import (
	"regexp"
	"unicode"
)

type PasswordUtility struct {
}

func CreatePasswordUtilty() *PasswordUtility {
	return &PasswordUtility{}
}

func (utility *PasswordUtility) IsStrong(plainPassword string) bool {
	if len(plainPassword) < 8 {
		return false
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

	return hasLower && hasUpper && hasDigit && hasSpecial
}
