package jwt

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func NewUtilMock() Handler {
	return &UtilMock{}
}

func (u *UtilMock) Generate(payload JwtPayload) (jwtString string, err error) {
	args := u.Called(payload)
	return args.String(0), args.Error(1)
}

func (u *UtilMock) Verify(jwtStr string) (valid bool, payload JwtPayload, err error) {
	args := u.Called(jwtStr)
	return args.Bool(0), args.Get(1).(JwtPayload), args.Error(3)
}
