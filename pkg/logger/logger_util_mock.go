package logger

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UtilMock struct {
	mock.Mock
}

func (u *UtilMock) Debug(msg string) {
	u.Called(msg)
}

func (u *UtilMock) Error(msg string) {
	u.Called(msg)
}

func (u *UtilMock) DebugCtx(ctx context.Context, msg string) {
	u.Called(ctx, msg)
}

func (u *UtilMock) ErrorCtx(ctx context.Context, msg string) {
	u.Called(ctx, msg)
}
