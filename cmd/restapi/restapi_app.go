package restapi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/nats"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

func StartApp() {
	appConfig := config.GetAppConfig("")

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
	validationHandler, err := validation.NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	// initialize core services
	jwtHandler, err := jwt.NewHandler(logService, appConfig)
	if err != nil {
		log.Fatal(err)
	}
	userService := user.NewService(logService, db)
	authService := auth.NewService(logService, jwtHandler, db)
	natsService, err := nats.NewPubService(appConfig, logService)
	if err != nil {
		log.Fatal(err)
	}
	defer natsService.Close()

	// initialize facades
	userFacade := user.NewFacade(appConfig, logService, userService, validationHandler, natsService)
	authFacade := auth.NewFacade(logService, authService, validationHandler)

	// register REST API routes
	router := RegisterRoutes(logService, authFacade, userFacade)

	// start HTTP server
	port := appConfig.APP_PORT
	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: router}
	go func() {
		logService.Debug(fmt.Sprintf("🚀 starting GO HTTP server on port: %s", port))
		err := server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logService.Debug("HTTP server closed")
			} else {
				logService.Debug(fmt.Sprintf("error while starting HTTP server: %v", err))
			}
		}
	}()

	// start NATS consumer
	go func() {
		err := natsService.Subscribe(appConfig.NATS_STREAM)
		if err != nil {
			logService.Error(err.Error())
		}
	}()

	// stop http server gracefully
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	err = server.Close()
	if err != nil {
		log.Fatal(fmt.Errorf("error closing HTTP server: %w", err))
	}
}
