package jwt

import "github.com/stretchr/testify/mock"

type UtilMockImpl struct {
	mock.Mock
}

func NewUtilMock() Util {
	return &UtilMockImpl{}
}

func (u *UtilMockImpl) Generate(userId string, userEmail string) (jwtString string, err error) {
	args := u.Called(userId, userEmail)
	return args.String(0), args.Error(1)
}
