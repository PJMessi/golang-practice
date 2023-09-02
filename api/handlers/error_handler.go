package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Type    string  `json:"type"`
	Message string  `json:"message"`
	Details *string `json:"details"`
}

// GlobalErrorHandler executes the handler function and returns 500 error response in case of panic
func GlobalErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic: %v", r)

				prepareInternalServerError(w)
			}
		}()

		next(w, r)
	}
}

// prepareErrorResponse updates the response writer with the provided status code, error type and error message
func prepareErrorResponse(w http.ResponseWriter, statusCode int, errorType string, errorMessage string, details *string) {
	w.Header().Set("Content-Type", "application/json")
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
func prepareInvalidRequestDataResponse(w http.ResponseWriter, details *string) {
	prepareErrorResponse(w, http.StatusBadRequest, "REQUEST_DATA.INVALID", "the provided data is invalid", details)
}

// prepareInternalServerError updates the response writer for an internal server error
func prepareInternalServerError(w http.ResponseWriter) {
	prepareErrorResponse(w, http.StatusInternalServerError, "INTERNAL", "internal server error", nil)
}

// FormatValidationErrors returns formated string describing the validation errors
func formatValidationErrors(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorMsg := ""
		for _, vErr := range errs {
			field := vErr.StructField()
			tag := vErr.Tag()
			errorMsg += fmt.Sprintf("'%s' validation failed for tag '%s'. ", field, tag)
		}
		return errorMsg
	}
	return err.Error()
}
