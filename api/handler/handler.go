package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pjmessi/go-database-usage/internal/business"
	"github.com/pjmessi/go-database-usage/internal/exceptions"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

type RoutesHandler struct {
	validator                  *validation.Validator
	accountRegistrationService *business.AccountRegistrationService
}

func CreateRouteHandler(
	validator *validation.Validator,
	accountRegistrationService *business.AccountRegistrationService,
) *RoutesHandler {
	return &RoutesHandler{
		validator:                  validator,
		accountRegistrationService: accountRegistrationService,
	}
}

type ErrorResponse struct {
	Type    string  `json:"type"`
	Message string  `json:"message"`
	Details *string `json:"details"`
}

func (routerHandler *RoutesHandler) HandleError(w http.ResponseWriter, err error) {
	log.Println("hell")
	log.Println(err)
	switch err := err.(type) {
	case *exceptions.InvalidRequestException:
		routerHandler.prepareErrorResponse(w, http.StatusBadRequest, err.Type, err.Message, nil)
	case *exceptions.DuplicateException:
		routerHandler.prepareErrorResponse(w, http.StatusBadRequest, err.Type, err.Message, nil)
	case *exceptions.NotFoundException:
		routerHandler.prepareErrorResponse(w, http.StatusNotFound, err.Type, err.Message, nil)
	default:
		routerHandler.prepareInternalServerError(w)
	}
}

// PanicHandler executes the handler function and returns 500 error response in case of panic
func (routerHandler *RoutesHandler) PanicHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic: %v", r)

				routerHandler.prepareInternalServerError(w)
			}
		}()

		next(w, r)
	}
}

// prepareErrorResponse updates the response writer with the provided status code, error type and error message
func (routerHandler *RoutesHandler) prepareErrorResponse(w http.ResponseWriter, statusCode int, errorType string, errorMessage string, details *string) {
	res, err := json.Marshal(ErrorResponse{Type: errorType, Message: errorMessage, Details: details})

	if err != nil {
		log.Printf("error while preparing error response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(res)
}

// prepareInvalidRequestDataResponse updates the response writer for a invalid request data error
func (routerHandler *RoutesHandler) prepareInvalidRequestDataResponse(w http.ResponseWriter, details *string) {
	routerHandler.prepareErrorResponse(w, http.StatusBadRequest, "REQUEST_DATA.INVALID", "the provided data is invalid", details)
}

// prepareInternalServerError updates the response writer for an internal server error
func (routerHandler *RoutesHandler) prepareInternalServerError(w http.ResponseWriter) {
	routerHandler.prepareErrorResponse(w, http.StatusInternalServerError, "INTERNAL", "internal server error", nil)
}
