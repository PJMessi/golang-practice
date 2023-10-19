package restapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/event"
	"github.com/pjmessi/golang-practice/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	// initialize database connection
	db, err := database.NewDb(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = db.CheckHealth()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.CloseConnection()

	// initialize common services
	logService := logger.NewService()
	validationUtil, err := validation.NewUtil()
	if err != nil {
		log.Fatal(err)
	}
	passwordUtil := password.NewUtil()
	jwtHandler, err := jwt.NewHandler(logService, appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// initialize core services
	userService := user.NewService(logService, db, passwordUtil)
	authService := auth.NewService(logService, jwtHandler, db, passwordUtil)
	eventPubService, err := event.NewPubService(appConfig, logService)
	if err != nil {
		log.Fatal(err)
	}

	// initialize facades
	userFacade := user.NewFacade(logService, userService, validationUtil, eventPubService)
	authFacade := auth.NewFacade(logService, authService, validationUtil)

	// register REST API routes
	router := RegisterRoutes(logService, authFacade, userFacade)

	// start http server
	port := appConfig.APP_PORT
	logService.Debug(fmt.Sprintf("🚀 starting GO server on port: %s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		logService.Debug(fmt.Sprintf("error while starting http server: %v", err))
	}
}
