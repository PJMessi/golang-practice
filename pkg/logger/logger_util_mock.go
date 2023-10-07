package logger

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UtilMockImpl struct {
	mock.Mock
}

func (u *UtilMockImpl) Debug(msg string) {
	u.Called(msg)
}

func (u *UtilMockImpl) Error(msg string) {
	u.Called(msg)
}

func (u *UtilMockImpl) DebugCtx(ctx context.Context, msg string) {
	u.Called(ctx, msg)
}

func (u *UtilMockImpl) ErrorCtx(ctx context.Context, msg string) {
	u.Called(ctx, msg)
}
