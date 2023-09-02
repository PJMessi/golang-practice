package hashing

import "golang.org/x/crypto/bcrypt"

type HashUtility struct {
}

func CreateHashUtility() *HashUtility {
	return &HashUtility{}
}

func (utility *HashUtility) HashString(plainString string) (*string, error) {
	hashValue, err := bcrypt.GenerateFromPassword([]byte(plainString), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashValueString := string(hashValue)

	return &hashValueString, nil
}

func (utility *HashUtility) VerifyHash(hashString string, plainString string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashString), []byte(plainString))
}
