package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pjmessi/go-database-usage/api"
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
	dbmysql "github.com/pjmessi/go-database-usage/internal/pkg/db/db-mysql"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	var dbInstance db.Db = dbmysql.CreateDbMysql()
	dbInstance.InitializeConnection(appConfig)
	defer dbInstance.CloseConnection()

	router := api.RegisterRoutes()

	appPort := fmt.Sprintf(":%s", appConfig.APP_PORT)

	log.Printf("starting server on port %s", appPort)
	if err := http.ListenAndServe(appPort, router); err != nil {
		log.Fatalf("error while starting http server: %v", err)
	}
}
