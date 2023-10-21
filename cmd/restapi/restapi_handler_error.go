package restapi

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/ctxutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/structutil"
	"github.com/pjmessi/golang-practice/pkg/strutil"
	"github.com/pjmessi/golang-practice/pkg/uuidutil"
)

type HttpHandlerWithCtx func(context.Context, http.ResponseWriter, *http.Request)
type HttpHandlerWithCtxAndJwtPayload func(context.Context, jwt.JwtPayload, http.ResponseWriter, *http.Request)

type ErrRes exception.Base

func (rh *RouteHandler) handleErr(ctx context.Context, w http.ResponseWriter, err error) {
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

func (rh *RouteHandler) handlePanic(next HttpHandlerWithCtx) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceId, err := uuidutil.GenUuidV4()
		if err != nil {
			rh.handleErr(context.Background(), w, fmt.Errorf("error while generating traceId"))
			return
		}
		ctx := ctxutil.NewCtxWithTraceId(traceId)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Trace-ID", traceId)

		defer func() {
			if recoverRes := recover(); recoverRes != nil {
				stack := make([]byte, 1024)
				runtime.Stack(stack, false)
				rh.logService.ErrorCtx(ctx, fmt.Sprintf("recovered from panic: %v\n%s", recoverRes, stack))
				rh.writeInternalErrRes(ctx, w)
			}
		}()

		startTime := time.Now()
		rh.logService.DebugCtx(ctx, fmt.Sprintf("new request: %s %s", r.Method, r.URL.String()))

		next(ctx, w, r)

		difference := time.Since(startTime)
		rh.logService.DebugCtx(ctx, fmt.Sprintf("request completed, took %d ms", difference.Milliseconds()))
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

func (rh *RouteHandler) handlePanicWithAuth(next HttpHandlerWithCtxAndJwtPayload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceId, err := uuidutil.GenUuidV4()
		if err != nil {
			rh.handleErr(context.Background(), w, fmt.Errorf("error while generating traceId"))
			return
		}
		ctx := ctxutil.NewCtxWithTraceId(traceId)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Trace-ID", traceId)

		defer func() {
			if recoverRes := recover(); recoverRes != nil {
				stack := make([]byte, 1024)
				runtime.Stack(stack, false)
				rh.logService.ErrorCtx(ctx, fmt.Sprintf("recovered from panic: %v\n%s", recoverRes, stack))
				rh.writeInternalErrRes(ctx, w)
			}
		}()

		jwt := rh.extractBearerToken(ctx, r)
		jwtPayload, err := rh.authFacade.VerifyJwt(ctx, jwt)
		if err != nil {
			rh.handleErr(ctx, w, err)
			return
		}

		startTime := time.Now()
		rh.logService.DebugCtx(ctx, fmt.Sprintf("new request: %s %s", r.Method, r.URL.String()))

		next(ctx, jwtPayload, w, r)

		difference := time.Since(startTime)
		rh.logService.DebugCtx(ctx, fmt.Sprintf("request completed, took %d ms", difference.Milliseconds()))
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
