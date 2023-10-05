package restapi

import (
	"context"
	"fmt"
	"log"
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
		rh.writeErrRes(w, http.StatusUnprocessableEntity, ErrRes(*e.Base))
	case exception.NotFound:
		rh.writeErrRes(w, http.StatusNotFound, ErrRes(*e.Base))
	case exception.Unauthenticated:
		rh.writeErrRes(w, http.StatusUnauthorized, ErrRes(*e.Base))
	case exception.Unauthorized:
		rh.writeErrRes(w, http.StatusForbidden, ErrRes(*e.Base))
	case exception.AlreadyExists:
		rh.writeErrRes(w, http.StatusBadRequest, ErrRes(*e.Base))
	case exception.FailedPrecondition:
		rh.writeErrRes(w, http.StatusBadRequest, ErrRes(*e.Base))
	default:
		log.Printf("unexpected error: %s", err.Error())
		rh.writeInternalErrRes(w)
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
				rh.writeInternalErrRes(w)
			}
		}()

		next(ctx, w, r)
	}
}

func (rh *RouteHandler) writeInternalErrRes(w http.ResponseWriter) {
	rh.writeErrRes(w, http.StatusInternalServerError, ErrRes{
		Type:    "INTERNAL",
		Message: "internal server error",
		Details: nil,
	})
}

func (rh *RouteHandler) writeErrRes(w http.ResponseWriter, statusCode int, errRes ErrRes) {
	resBytes, err := structutil.ConvertToBytes(errRes)
	if err != nil {
		log.Printf("err while converting ErrRes to bytes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Printf("err while writing err response: %v\n", err)
		}
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resBytes)
	if err != nil {
		log.Printf("err while writing err response: %v\n", err)
	}
}
