package app

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var dbInstance db.Db = db.CreateDbMysql()
	dbInstance.InitializeConnection(appConfig)
	defer dbInstance.CloseConnection()
}
