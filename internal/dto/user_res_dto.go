package dto

import (
	"time"

	"github.com/pjmessi/golang-practice/internal/model"
)

func UserToUserRes(user *model.User) model.UserRes {
	return model.UserRes{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
