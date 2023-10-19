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
