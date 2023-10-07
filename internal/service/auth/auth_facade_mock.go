package auth

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type FacadeMock struct {
	mock.Mock
}

func (f *FacadeMock) Login(ctx context.Context, reqBytes []byte) ([]byte, error) {
	args := f.Called(ctx, reqBytes)
	return args.Get(0).([]byte), args.Error(1)
}
