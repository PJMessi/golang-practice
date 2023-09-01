package app

import (
	"fmt"
	"log"

	godotenv "github.com/joho/godotenv"
	"github.com/pjmessi/go-database-usage/config"
	db "github.com/pjmessi/go-database-usage/internal/pkg/db"
)

func StartApp() {
	fmt.Println("StartApp >>>")

	appConfig := config.GetAppConfig()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var dbInstance db.Db = db.CreateDbMysql()
	dbInstance.InitializeConnection(appConfig)
	defer dbInstance.CloseConnection()

	fmt.Println("<<< StartApp")
}
