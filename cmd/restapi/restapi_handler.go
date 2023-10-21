package restapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/ctxutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/structutil"
	"github.com/pjmessi/golang-practice/pkg/strutil"
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

func (rh *RouteHandler) attachMiddlewares(handlerFunc http.HandlerFunc, authenticate bool) http.HandlerFunc {
	handler := http.Handler(handlerFunc)

	if authenticate {
		handler = rh.authMiddleware(handler)
	}

	handler = rh.reqLoggerMiddleware(handler)
	handler = rh.panicHandlerMiddleware(handler)
	handler = rh.ctxMiddleware(handler)
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
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

func (rh *RouteHandler) writeHttpResFromErr(ctx context.Context, w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case exception.InvalidReq:
		rh.convertDetailsKeyToCamelcase(&e)
		rh.writeErrRes(ctx, w, http.StatusUnprocessableEntity, ErrRes(*e.Base))
	case exception.NotFound:
		rh.writeErrRes(ctx, w, http.StatusNotFound, ErrRes(*e.Base))
	case exception.Unauthenticated:
		rh.writeErrRes(ctx, w, http.StatusUnauthorized, ErrRes(*e.Base))
	case exception.Unauthorized:
		rh.writeErrRes(ctx, w, http.StatusForbidden, ErrRes(*e.Base))
	case exception.AlreadyExists:
		rh.writeErrRes(ctx, w, http.StatusBadRequest, ErrRes(*e.Base))
	case exception.FailedPrecondition:
		rh.writeErrRes(ctx, w, http.StatusBadRequest, ErrRes(*e.Base))
	default:
		rh.logService.ErrorCtx(ctx, fmt.Sprintf("unexpected error: %s", err.Error()))
		rh.writeInternalErrRes(ctx, w)
	}
}

func (rh *RouteHandler) writeInternalErrRes(ctx context.Context, w http.ResponseWriter) {
	rh.writeErrRes(ctx, w, http.StatusInternalServerError, ErrRes{
		Type:    "INTERNAL",
		Message: "internal server error",
		Details: nil,
	})
}

func (rh *RouteHandler) writeErrRes(ctx context.Context, w http.ResponseWriter, statusCode int, errRes ErrRes) {
	resBytes, err := structutil.ConvertToBytes(errRes)
	if err != nil {
		rh.logService.ErrorCtx(ctx, fmt.Sprintf("err while converting ErrRes to bytes: %v", err))
		w.WriteHeader(http.StatusInternalServerError)

		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			rh.logService.ErrorCtx(ctx, fmt.Sprintf("err while writing err response: %v", err))
		}
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resBytes)
	if err != nil {
		rh.logService.ErrorCtx(ctx, fmt.Sprintf("err while writing err response: %v", err))
	}
}

func (rh *RouteHandler) convertDetailsKeyToCamelcase(validationErr *exception.InvalidReq) {
	camelCaseDetails := map[string]string{}
	if validationErr.Details != nil {
		for key, val := range *validationErr.Details {
			camelcaseKey := strutil.PascalCaseToCamelCase(key)
			camelCaseDetails[camelcaseKey] = val
		}
		*validationErr.Details = camelCaseDetails
	}
}

func (rh *RouteHandler) extractBearerToken(ctx context.Context, r *http.Request) string {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		rh.logService.DebugCtx(ctx, "empty Authorization header")
		return ""
	}

	jwt, prefixExists := strings.CutPrefix(authHeader, "Bearer ")
	if !prefixExists {
		rh.logService.DebugCtx(ctx, "token does not start with 'Bearer '")
		return ""
	}

	return jwt
}

func (rh *RouteHandler) handleRouteNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := exception.NewNotFoundFromBase(exception.Base{
			Type:    "ROUTE.NOT_FOUND",
			Message: "route not found",
		})

		rh.writeHttpResFromErr(ctx, w, err)
	}
}
