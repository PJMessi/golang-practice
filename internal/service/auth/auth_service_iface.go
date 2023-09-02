package auth

import (
	"github.com/pjmessi/go-database-usage/internal/model"
)

type Service interface {
	Login(email string, password string) (user model.User, jwt string, err error)
}
