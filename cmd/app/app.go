package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pjmessi/go-database-usage/api"
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
	dbmysql "github.com/pjmessi/go-database-usage/internal/pkg/db/db-mysql"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	// initialize validator
	validator := validation.CreateValidator()
	validator.InitializeValidator()

	// initialize database connection
	var dbInstance db.Db = dbmysql.CreateDbMysql()
	dbInstance.InitializeConnection(appConfig)
	defer dbInstance.CloseConnection()

	// register REST API routes
	router := api.RegisterRoutes(validator)

	// start http server
	appPort := fmt.Sprintf(":%s", appConfig.APP_PORT)
	log.Printf("starting server on port %s", appPort)
	if err := http.ListenAndServe(appPort, router); err != nil {
		log.Fatalf("error while starting http server: %v", err)
	}
}
