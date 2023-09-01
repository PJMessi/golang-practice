package app

import (
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	var dbInstance db.Db = db.CreateDbMysql()
	dbInstance.InitializeConnection(appConfig)
	defer dbInstance.CloseConnection()
}
