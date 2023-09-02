package dtos

import (
	"github.com/pjmessi/go-database-usage/api/responses"
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
)

func UserToUserResponse(user *model.User) *responses.UserResponse {
	return &responses.UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}
