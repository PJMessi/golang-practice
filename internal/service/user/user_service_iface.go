package user

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string) (model.User, error)
}
