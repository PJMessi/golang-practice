package restapi

import (
	"context"
	"net/http"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"

	"github.com/gorilla/mux"
)

type FacadeApiFunc func(ctx context.Context, reqBytes []byte) ([]byte, error)
type FacadeApiFuncWithAuth func(ctx context.Context, reqBytes []byte, jwtPayload jwt.JwtPayload) ([]byte, error)
type ErrRes exception.Base

func RegisterRoutes(logService logger.Service, authFacade auth.Facade, userFacade user.Facade) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	routeHandler := NewRouteHandler(logService, authFacade, userFacade)

	router.HandleFunc("/auth/login", routeHandler.attachMiddlewares(routeHandler.handlePublicApi(authFacade.Login), false)).Methods("POST")

	router.HandleFunc("/users/registration", routeHandler.attachMiddlewares(routeHandler.handlePublicApi(userFacade.RegisterUser), false)).Methods("POST")
	router.HandleFunc("/users/profile", routeHandler.attachMiddlewares(routeHandler.handlePrivateApi(userFacade.GetProfile), true)).Methods("GET")

	router.NotFoundHandler = routeHandler.attachMiddlewares(routeHandler.handleRouteNotFound(), false)
	return router
}
