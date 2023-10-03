package database

import (
	"context"

	"github.com/pjmessi/go-database-usage/internal/model"
)

type Db interface {
	CloseConnection()
	IsHealthy() bool

	CreateUser(ctx context.Context, user *model.User) error
	IsUserEmailTaken(ctx context.Context, email string) (isTaken bool, err error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user model.User, err error)
}
