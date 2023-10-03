package user

import (
	"context"

	"github.com/pjmessi/go-database-usage/internal/model"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string) (model.User, error)
}
