package jwt

import "github.com/stretchr/testify/mock"

type HandlerMock struct {
	mock.Mock
}

func NewUtilMock() Handler {
	return &HandlerMock{}
}

func (h *HandlerMock) Generate(payload JwtPayload) (jwtString string, err error) {
	args := h.Called(payload)
	return args.String(0), args.Error(1)
}

func (h *HandlerMock) Verify(jwtStr string) (valid bool, payload JwtPayload, err error) {
	args := h.Called(jwtStr)
	return args.Bool(0), args.Get(1).(JwtPayload), args.Error(3)
}
