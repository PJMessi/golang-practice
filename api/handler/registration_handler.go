package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pjmessi/go-database-usage/api/responses"
	"github.com/pjmessi/go-database-usage/internal/dtos"
)

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=255"`
}

type RegisterUserResponse struct {
	User *responses.UserResponse `json:"user"`
}

func (routesHandler *RoutesHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest RegisterUserRequest
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

	user, err := routesHandler.accountRegistrationService.RegisterUser(
		registerUserRequest.Email,
		registerUserRequest.Password,
	)
	if err != nil {
		routesHandler.HandleError(w, err)
		return
	}

	response := &RegisterUserResponse{
		User: dtos.UserToUserResponse(user),
	}

	responseInBytes, err := json.Marshal(response)
	if err != nil {
		routesHandler.HandleError(w, err)
		return
	}

	w.Write(responseInBytes)
}
