package logger

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Debug(msg string) {
	s.Called(msg)
}

func (s *ServiceMock) Error(msg string) {
	s.Called(msg)
}

func (s *ServiceMock) DebugCtx(ctx context.Context, msg string) {
	s.Called(ctx, msg)
}

func (s *ServiceMock) ErrorCtx(ctx context.Context, msg string) {
	s.Called(ctx, msg)
}
