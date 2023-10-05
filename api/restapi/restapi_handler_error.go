package restapi

import (
	"log"
	"net/http"
	"runtime"

	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/structutil"
)

type ErrRes exception.Base

func (rh *RouteHandler) handleErr(w http.ResponseWriter, err error) {
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

func (rh *RouteHandler) handlePanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer func() {
			if recoverRes := recover(); recoverRes != nil {
				stack := make([]byte, 1024)
				runtime.Stack(stack, false)

				log.Printf("recovered from panic: %v\n%s", recoverRes, stack)

				rh.writeInternalErrRes(w)
			}
		}()

		next(w, r)
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
