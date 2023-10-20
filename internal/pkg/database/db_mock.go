package database

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (r *DbMock) CheckHealth() error {
	args := r.Called()
	return args.Error(0)
}

func (r *DbMock) SaveUser(ctx context.Context, user *model.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *DbMock) IsUserEmailTaken(ctx context.Context, email string) (bool, error) {
	args := r.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (r *DbMock) GetUserByEmail(ctx context.Context, email string) (bool, model.User, error) {
	args := r.Called(ctx, email)
	return args.Bool(0), args.Get(1).(model.User), args.Error(2)
}

func (r *DbMock) GetUserById(ctx context.Context, userId string) (bool, model.User, error) {
	args := r.Called(ctx, userId)
	return args.Bool(0), args.Get(1).(model.User), args.Error(2)
}

func (r *DbMock) CloseConnection() {
	r.Called()
}
