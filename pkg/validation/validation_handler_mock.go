package validation

import "github.com/stretchr/testify/mock"

type HandlerMock struct {
	mock.Mock
}

func (u *HandlerMock) ValidateStruct(s interface{}) error {
	args := u.Called(s)
	return args.Error(0)
}

func (u *HandlerMock) FormatValidationError(err error) string {
	args := u.Called(err)
	return args.String(0)
}
