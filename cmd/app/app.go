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
	validationUtil := validation.NewUtil()
	hashUtil := hash.NewUtil()
	passwordUtil := password.NewUtil()
	uuidUtil := uuid.NewUtil()
	jwtUtility, err := jwt.NewUtil(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// initialize services
	userService := user.NewService(db, hashUtil, passwordUtil, uuidUtil)
	authService := auth.NewService(jwtUtility, db, hashUtil)

	// initialize facades
	userFacade := user.NewFacade(userService, validationUtil)
	authFacade := auth.NewFacade(authService, validationUtil)

	// register REST API routes
	router := restapi.RegisterRoutes(authFacade, userFacade, uuidUtil)

	// start http server
	appPort := fmt.Sprintf(":%s", appConfig.APP_PORT)
	log.Printf("starting server on port %s", appPort)
	if err := http.ListenAndServe(appPort, router); err != nil {
		log.Fatalf("error while starting http server: %v", err)
	}
}
