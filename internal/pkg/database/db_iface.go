package database

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
)

type Db interface {
	CloseConnection()
	CheckHealth() error

	SaveUser(ctx context.Context, user *model.User) error
	IsUserEmailTaken(ctx context.Context, email string) (isTaken bool, err error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user model.User, err error)
	GetUserById(ctx context.Context, userId string) (exists bool, user model.User, err error)
}
