package auth

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (user model.User, jwt string, err error)
}
