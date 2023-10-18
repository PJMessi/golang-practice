package restapi

import (
	"net/http"

	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"

	"github.com/gorilla/mux"
)

func RegisterRoutes(loggerUtil logger.Util, authFacade auth.Facade, userFacade user.Facade) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	handler := NewRouteHandler(loggerUtil, authFacade, userFacade)

	router.HandleFunc("/auth/login", handler.handlePanic(handler.handleLoginApi)).Methods("POST")
	router.HandleFunc("/user/registration", handler.handlePanic(handler.handleUserRegApi)).Methods("POST")

	return router
}
