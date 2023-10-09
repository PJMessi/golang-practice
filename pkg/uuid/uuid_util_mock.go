package uuid

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func (u *UtilMock) GenUuidV4() (string, error) {
	args := u.Called()
	return args.String(0), args.Error(1)
}
