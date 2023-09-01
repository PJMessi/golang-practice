package db

import (
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
)

type Db interface {
	InitializeConnection(appConfig *config.AppConfig)
	CloseConnection()
	CreateUser() *model.User
	GetUser() *model.User
	IsHealthy() bool
}
