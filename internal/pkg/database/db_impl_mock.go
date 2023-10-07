package database

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/stretchr/testify/mock"
)

type DbMockImpl struct {
	mock.Mock
}

func (r *DbMockImpl) IsHealthy() bool {
	args := r.Called()
	return args.Bool(0)
}

func (r *DbMockImpl) CreateUser(ctx context.Context, user *model.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *DbMockImpl) IsUserEmailTaken(ctx context.Context, email string) (bool, error) {
	args := r.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (r *DbMockImpl) GetUserByEmail(ctx context.Context, email string) (bool, model.User, error) {
	args := r.Called(ctx, email)
	return args.Bool(0), args.Get(1).(model.User), args.Error(2)
}

func (r *DbMockImpl) CloseConnection() {
	r.Called()
}
