package restapi

import (
	"context"
	"io"
	"net/http"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/ctxutil"
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

	return router
}

func (rh *RouteHandler) handlePublicApi(facadeFunc FacadeApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}

		// TODO: TEST
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				rh.writeHttpResFromErr(ctx, w, err)
			}
		}(r.Body)

		resByte, err := facadeFunc(ctx, reqBytes)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}

		_, err = w.Write(resByte)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}
	}

}

func (rh *RouteHandler) handlePrivateApi(facadeFunc FacadeApiFuncWithAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		jwtPayload, ok := ctxutil.GetValue(ctx, "jwtPayload").(jwt.JwtPayload)
		if !ok {
			rh.logService.DebugCtx(ctx, "restapi.RouteHandler.handlePrivateApi(): jwtPayload not set in context")
			rh.writeHttpResFromErr(ctx, w, exception.NewUnauthenticated())
			return
		}

		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}

		// TODO: TEST
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				rh.writeHttpResFromErr(ctx, w, err)
			}
		}(r.Body)

		resByte, err := facadeFunc(ctx, reqBytes, jwtPayload)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}

		_, err = w.Write(resByte)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}
	}

}
