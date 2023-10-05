package restapi

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/pjmessi/golang-practice/pkg/ctxutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/structutil"
)

type HttpHandlerWithCtx func(context.Context, http.ResponseWriter, *http.Request)

type ErrRes exception.Base

func (rh *RouteHandler) handleErr(ctx context.Context, w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case exception.InvalidReq:
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
		rh.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("unexpected error: %s", err.Error()))
		rh.writeInternalErrRes(ctx, w)
	}
}

func (rh *RouteHandler) handlePanic(next HttpHandlerWithCtx) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceId, err := rh.uuidUtil.GenUuidV4()
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
				rh.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("recovered from panice: %v\n%s", recoverRes, stack))
				rh.writeInternalErrRes(ctx, w)
			}
		}()

		next(ctx, w, r)
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
		rh.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("err while converting ErrRes to bytes: %v", err))
		w.WriteHeader(http.StatusInternalServerError)

		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			rh.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("err while writing err response: %v", err))
		}
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resBytes)
	if err != nil {
		rh.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("err while writing err response: %v", err))
	}
}
