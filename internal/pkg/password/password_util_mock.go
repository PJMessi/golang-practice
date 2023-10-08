package password

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func (u *UtilMock) IsStrong(plainPw string) bool {
	args := u.Called(plainPw)
	return args.Bool(0)
}

func (u *UtilMock) Hash(plainPw string) (string, error) {
	args := u.Called(plainPw)
	return args.String(0), args.Error(1)
}

func (u *UtilMock) IsHashCorrect(hashedPw string, plainPw string) bool {
	args := u.Called(hashedPw, plainPw)
	return args.Bool(0)
}
