package auth

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Login(ctx context.Context, email string, password string) (model.User, string, error) {
	args := s.Called(ctx, email, password)
	return args.Get(0).(model.User), args.String(1), args.Error(2)
}

func (s *ServiceMock) VerifyJwt(ctx context.Context, jwtStr string) (jwt.JwtPayload, error) {
	args := s.Called(ctx, jwtStr)
	return args.Get(0).(jwt.JwtPayload), args.Error(1)
}
