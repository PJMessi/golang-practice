package restapi

import (
	"net/http"

	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"

	"github.com/gorilla/mux"
)

func RegisterRoutes(logService logger.Service, authFacade auth.Facade, userFacade user.Facade) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	handler := NewRouteHandler(logService, authFacade, userFacade)

	router.HandleFunc("/auth/login", handler.handlePanic(handler.handleLoginApi)).Methods("POST")
	router.HandleFunc("/users/registration", handler.handlePanic(handler.handleUserRegApi)).Methods("POST")

	return router
}
