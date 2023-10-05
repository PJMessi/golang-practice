package hash

import (
	"github.com/pjmessi/golang-practice/pkg/strutil"
	"golang.org/x/crypto/bcrypt"
)

type UtilImpl struct {
	Util
}

func NewUtil() Util {
	return &UtilImpl{}
}

func (utility *UtilImpl) GenerateHash(plainStr string) (hashStr string, err error) {
	plainStrBytes := strutil.ConvertToBytes(plainStr)
	hashedStrBytes, err := bcrypt.GenerateFromPassword(plainStrBytes, bcrypt.DefaultCost)
	if err != nil {
		return hashStr, err
	}
	hashStr = strutil.ConvertFromBytes(hashedStrBytes)
	return hashStr, nil
}

func (utility *UtilImpl) VerifyHash(hashStr string, plainStr string) bool {
	hashStrBytes := strutil.ConvertToBytes(hashStr)
	plainStrBytes := strutil.ConvertToBytes(plainStr)
	err := bcrypt.CompareHashAndPassword(hashStrBytes, plainStrBytes)
	return err == nil
}
