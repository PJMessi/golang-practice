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
	routeHandler := NewRouteHandler(logService, authFacade, userFacade)

	router.HandleFunc("/auth/login", routeHandler.attachMiddlewares(routeHandler.handleLoginApi, false)).Methods("POST")
	router.HandleFunc("/users/registration", routeHandler.attachMiddlewares(routeHandler.handleUserRegApi, false)).Methods("POST")
	router.HandleFunc("/users/profile", routeHandler.attachMiddlewares(routeHandler.handleGetProfileApi, true)).Methods("GET")

	return router
}
