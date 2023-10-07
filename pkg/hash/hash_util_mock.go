package hash

import "github.com/stretchr/testify/mock"

type UtilMockImpl struct {
	mock.Mock
}

func (u *UtilMockImpl) GenerateHash(plainString string) (hashedString string, err error) {
	args := u.Called(plainString)
	return args.String(0), args.Error(1)
}

func (u *UtilMockImpl) VerifyHash(hashString string, plainString string) bool {
	args := u.Called(hashString, plainString)
	return args.Bool(0)
}
