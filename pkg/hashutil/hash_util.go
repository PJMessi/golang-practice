package hashutil

import (
	"fmt"

	"github.com/pjmessi/golang-practice/pkg/strutil"
	"golang.org/x/crypto/bcrypt"
)

func Generate(plainStr string) (hashStr string, err error) {
	plainStrBytes := strutil.ConvertToBytes(plainStr)
	hashedStrBytes, err := bcrypt.GenerateFromPassword(plainStrBytes, bcrypt.DefaultCost)
	if err != nil {
		return hashStr, fmt.Errorf("hashutil.GenerateHash(): %w", err)
	}
	hashStr = strutil.ConvertFromBytes(hashedStrBytes)
	return hashStr, nil
}

func Verify(hashStr string, plainStr string) bool {
	hashStrBytes := strutil.ConvertToBytes(hashStr)
	plainStrBytes := strutil.ConvertToBytes(plainStr)
	err := bcrypt.CompareHashAndPassword(hashStrBytes, plainStrBytes)
	return err == nil
}
