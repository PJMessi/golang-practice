package db

import (
	"github.com/pjmessi/go-database-usage/config"
)

type Db interface {
	InitializeConnection(appConfig *config.AppConfig)
	CloseConnection()
	IsHealthy() bool
}
