package db

import (
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
)

type Db interface {
	InitializeConnection(appConfig *config.AppConfig)
	CloseConnection()
	IsHealthy() bool

	// User related operations
	CreateUser(user *model.User) error
	IsUserEmailTaken(email string) (*bool, error)
}
