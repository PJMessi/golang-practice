package restapi

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/ctxutil"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/uuidutil"
)

type RouteHandler struct {
	authFacade auth.Facade
	userFacade user.Facade
	logService logger.Service
}

func NewRouteHandler(logService logger.Service, authFacade auth.Facade, userFacade user.Facade) *RouteHandler {
	return &RouteHandler{
		authFacade: authFacade,
		userFacade: userFacade,
		logService: logService,
	}
}

func (rh *RouteHandler) attachMiddlewares(handler http.HandlerFunc, authenticate bool) http.HandlerFunc {
	apiHandler := http.Handler(handler)

	if authenticate {
		apiHandler = rh.authMiddleware(apiHandler)
	}

	apiHandler = rh.reqLoggerMiddleware(apiHandler)
	apiHandler = rh.panicHandlerMiddleware(apiHandler)
	apiHandler = rh.ctxMiddleware(apiHandler)
	return func(w http.ResponseWriter, r *http.Request) {
		apiHandler.ServeHTTP(w, r)
	}
}

func (rh *RouteHandler) ctxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId, err := uuidutil.GenUuidV4()
		if err != nil {
			rh.writeHttpResFromErr(context.Background(), w, fmt.Errorf("error while generating traceId"))
			return
		}
		ctx := ctxutil.NewCtxWithTraceId(traceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rh *RouteHandler) panicHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceId := ctxutil.GetTraceIdFromCtx(ctx)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Trace-ID", traceId)

		// prints the stack trace on panic and writes 500 http respnose
		defer func() {
			if recoverRes := recover(); recoverRes != nil {
				stack := make([]byte, 1024)
				runtime.Stack(stack, false)
				rh.logService.ErrorCtx(ctx, fmt.Sprintf("recovered from panic: %v\n%s", recoverRes, stack))
				rh.writeInternalErrRes(ctx, w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (rh *RouteHandler) reqLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		startTime := time.Now()
		rh.logService.DebugCtx(ctx, fmt.Sprintf("new request: %s %s", r.Method, r.URL.String()))

		next.ServeHTTP(w, r)

		difference := time.Since(startTime)
		rh.logService.DebugCtx(ctx, fmt.Sprintf("request completed, took %d ms", difference.Milliseconds()))
	})
}

func (rh *RouteHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		jwt := rh.extractBearerToken(ctx, r)
		jwtPayload, err := rh.authFacade.VerifyJwt(ctx, jwt)
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
			return
		}

		ctx = ctxutil.AddValue(ctx, "jwtPayload", jwtPayload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
