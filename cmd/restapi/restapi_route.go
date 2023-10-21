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
	rHandler := NewRouteHandler(logService, authFacade, userFacade)

	// auth routes
	router.HandleFunc("/auth/login", rHandler.attachMiddlewares(rHandler.handlePublicApi(authFacade.Login), false)).Methods("POST")

	// user routes
	router.HandleFunc("/users/registration", rHandler.attachMiddlewares(rHandler.handlePublicApi(userFacade.RegisterUser), false)).Methods("POST")
	router.HandleFunc("/users/profile", rHandler.attachMiddlewares(rHandler.handlePrivateApi(userFacade.GetProfile), true)).Methods("GET")

	router.NotFoundHandler = rHandler.attachMiddlewares(rHandler.handleRouteNotFound(), false)
	return router
}
