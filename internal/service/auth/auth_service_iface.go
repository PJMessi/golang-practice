package auth

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (user model.User, jwt string, err error)
	VerifyJwt(ctx context.Context, jwtStr string) (jwt.JwtPayload, error)
}
