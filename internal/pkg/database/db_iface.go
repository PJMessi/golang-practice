package database

import (
	"github.com/pjmessi/go-database-usage/internal/model"
)

type Db interface {
	CloseConnection()
	IsHealthy() bool

	CreateUser(user *model.User) error
	IsUserEmailTaken(email string) (isTaken bool, err error)
	GetUserByEmail(email string) (exists bool, user model.User, err error)
}
