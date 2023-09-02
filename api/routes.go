package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pjmessi/go-database-usage/api/handlers"
)

func RegisterRoutes() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/account/registration", handlers.GlobalErrorHandler(handlers.RegisterUserHandler)).Methods("POST")

	return router
}
