package handler

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (routesHandler *RoutesHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&registerUserRequest); err != nil {
		routesHandler.prepareInvalidRequestDataResponse(w, nil)
		return
	}

	err := routesHandler.validator.ValidateStruct(registerUserRequest)
	if err != nil {
		details := routesHandler.validator.FormatValidationError(err)
		routesHandler.prepareInvalidRequestDataResponse(w, &details)
		return
	}
}
