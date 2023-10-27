package tests

import (
	"log"
	"net/http/httptest"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pjmessi/golang-practice/cmd/restapi"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/nats"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

var testServer *httptest.Server
var db database.Db
var appConfig *config.AppConfig
var testDbCon *sql.DB

func setupIntegrationTest() {
	appConfig = config.GetAppConfig("test")

	var err error

	// initialize database connection
	db, err = database.NewDb(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = db.CheckHealth()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// initialize common services
	logService := logger.NewService()
	validationHandler, err := validation.NewHandler()
	if err != nil {
		log.Fatal(err)
	}
	jwtHandler, err := jwt.NewHandler(logService, appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// initialize core services
	userService := user.NewService(logService, db)
	authService := auth.NewService(logService, jwtHandler, db)
	natsService, err := nats.NewPubService(appConfig, logService)
	if err != nil {
		log.Fatal(err)
	}

	// initialize facades
	userFacade := user.NewFacade(appConfig, logService, userService, validationHandler, natsService)
	authFacade := auth.NewFacade(logService, authService, validationHandler)

	// register REST API routes
	router := restapi.RegisterRoutes(logService, authFacade, userFacade)

	// start http server
	testServer = httptest.NewServer(router)

	// initialize database connection for testing
	testDbCon, err = testutil.GetTestDbCon(appConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func teardownIntegrationTest() {
	testDbCon.Exec("DELETE FROM users;")

	// Clean up resources and shut down the test server and test database
	testDbCon.Close()
	db.CloseConnection()
	testServer.Close()
	// Additional cleanup as needed
}
