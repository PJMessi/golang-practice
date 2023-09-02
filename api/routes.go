package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pjmessi/go-database-usage/api/handler"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

func RegisterRoutes(validator *validation.Validator) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	routeHandler := handler.CreateRouteHandler(validator)

	router.HandleFunc("/account/registration", routeHandler.GlobalErrorHandler(routeHandler.RegisterUserHandler)).Methods("POST")

	return router
}
