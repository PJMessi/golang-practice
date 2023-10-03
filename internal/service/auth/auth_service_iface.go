package auth

import (
	"context"

	"github.com/pjmessi/go-database-usage/internal/model"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (user model.User, jwt string, err error)
}
