package validation

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func (u *UtilMock) ValidateStruct(s interface{}) error {
	args := u.Called(s)
	return args.Error(0)
}

func (u *UtilMock) FormatValidationError(err error) string {
	args := u.Called(err)
	return args.String(0)
}
