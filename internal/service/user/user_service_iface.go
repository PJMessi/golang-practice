package user

import "github.com/pjmessi/go-database-usage/internal/model"

type Service interface {
	CreateUser(email string, password string) (model.User, error)
}
