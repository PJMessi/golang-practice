package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pjmessi/golang-practice/api/restapi"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/hash"
	"github.com/pjmessi/golang-practice/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/password"
	"github.com/pjmessi/golang-practice/pkg/uuid"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

func StartApp() {
	appConfig := config.GetAppConfig()

	// initialize database connection
	var db, err = database.NewDbImpl(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection()

	// initialize utilities
	loggerUtil := logger.NewUtil()
	validationUtil := validation.NewUtil()
	hashUtil := hash.NewUtil()
	passwordUtil := password.NewUtil()
	uuidUtil := uuid.NewUtil()
	jwtUtility, err := jwt.NewUtil(loggerUtil, appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// initialize services
	userService := user.NewService(loggerUtil, db, hashUtil, passwordUtil, uuidUtil)
	authService := auth.NewService(loggerUtil, jwtUtility, db, hashUtil)

	// initialize facades
	userFacade := user.NewFacade(loggerUtil, userService, validationUtil)
	authFacade := auth.NewFacade(loggerUtil, authService, validationUtil)

	// register REST API routes
	router := restapi.RegisterRoutes(loggerUtil, authFacade, userFacade, uuidUtil)

	// start http server
	port := appConfig.APP_PORT
	loggerUtil.Debug(fmt.Sprintf("🚀 starting GO server on port: %s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		loggerUtil.Debug(fmt.Sprintf("error while starting http server: %v", err))
	}
}
