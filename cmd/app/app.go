package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pjmessi/go-database-usage/api/restapi"
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/database"
	"github.com/pjmessi/go-database-usage/internal/service/auth"
	"github.com/pjmessi/go-database-usage/internal/service/user"
	"github.com/pjmessi/go-database-usage/pkg/hash"
	"github.com/pjmessi/go-database-usage/pkg/jwt"
	"github.com/pjmessi/go-database-usage/pkg/logger"
	"github.com/pjmessi/go-database-usage/pkg/password"
	"github.com/pjmessi/go-database-usage/pkg/uuid"
	"github.com/pjmessi/go-database-usage/pkg/validation"
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
	jwtUtility, err := jwt.NewUtil(appConfig)
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
	appPort := fmt.Sprintf(":%s", appConfig.APP_PORT)
	loggerUtil.Debug(fmt.Sprintf("starting server on port %s", appPort))
	if err := http.ListenAndServe(appPort, router); err != nil {
		loggerUtil.Debug(fmt.Sprintf("error while starting http server: %v", err))
	}
}
