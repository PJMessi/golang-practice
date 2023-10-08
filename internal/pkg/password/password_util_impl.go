package password

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/pjmessi/golang-practice/pkg/hash"
)

type UtilImpl struct {
	hashUtil hash.Util
}

func NewUtil(hashUtil hash.Util) Util {
	return &UtilImpl{
		hashUtil: hashUtil,
	}
}

func (utility *UtilImpl) IsStrong(plainPw string) bool {
	// password should be at least 8 characters long and must have at least 1 lowercase character, 1 uppercase character,
	// 1 digit and 1 special character (!@#$%^&*()_+{})

	if len(plainPw) < 8 {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range plainPw {
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

func (utility *UtilImpl) Hash(plainPw string) (string, error) {
	hashedPw, err := utility.hashUtil.GenerateHash(plainPw)
	if err != nil {
		return "", fmt.Errorf("password.Hash: %w", err)
	}

	return hashedPw, nil
}

func (utility *UtilImpl) IsHashCorrect(hashedPw string, plainPw string) bool {
	isValidHash := utility.hashUtil.VerifyHash(hashedPw, plainPw)
	return isValidHash
}
