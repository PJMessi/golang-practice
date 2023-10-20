package jwt

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func NewUtilMock() Handler {
	return &UtilMock{}
}

func (u *UtilMock) Generate(userId string, userEmail string) (jwtString string, err error) {
	args := u.Called(userId, userEmail)
	return args.String(0), args.Error(1)
}

func (u *UtilMock) Verify(jwtStr string) (valid bool, userId string, userEmail string, err error) {
	args := u.Called(jwtStr)
	return args.Bool(0), args.String(1), args.String(2), args.Error(3)
}
