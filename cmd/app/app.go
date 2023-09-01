package app

import (
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
	dbmysql "github.com/pjmessi/go-database-usage/internal/pkg/db/db-mysql"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	var dbInstance db.Db = dbmysql.CreateDbMysql()
	dbInstance.InitializeConnection(appConfig)
	defer dbInstance.CloseConnection()
}
