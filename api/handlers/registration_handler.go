package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&registerUserRequest); err != nil {
		prepareInvalidRequestDataResponse(w, nil)
		return
	}

	validator := validator.New()
	err := validator.Struct(registerUserRequest)
	if err != nil {
		details := formatValidationErrors(err)
		prepareInvalidRequestDataResponse(w, &details)
		return
	}
}
