package restapi

import (
	"net/http"

	"github.com/pjmessi/go-database-usage/internal/service/auth"
	"github.com/pjmessi/go-database-usage/internal/service/user"

	"github.com/gorilla/mux"
)

func RegisterRoutes(authFacade auth.Facade, userFacade user.Facade) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	handler := NewRouteHandler(authFacade, userFacade)

	router.HandleFunc("/auth/login", handler.handlePanic(handler.handleLoginApi)).Methods("POST")
	router.HandleFunc("/user/registration", handler.handlePanic(handler.handleUserRegApi)).Methods("POST")

	return router
}
