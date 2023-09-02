package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pjmessi/go-database-usage/api/handler"
	"github.com/pjmessi/go-database-usage/internal/business"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

func RegisterRoutes(validator *validation.Validator, accountRegistrationService *business.AccountRegistrationService) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	routeHandler := handler.CreateRouteHandler(validator, accountRegistrationService)

	router.HandleFunc("/account/registration", routeHandler.PanicHandler(routeHandler.RegisterUserHandler)).Methods("POST")

	return router
}
