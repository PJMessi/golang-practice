package hash

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func (u *UtilMock) GenerateHash(plainString string) (hashedString string, err error) {
	args := u.Called(plainString)
	return args.String(0), args.Error(1)
}

func (u *UtilMock) VerifyHash(hashString string, plainString string) bool {
	args := u.Called(hashString, plainString)
	return args.Bool(0)
}
