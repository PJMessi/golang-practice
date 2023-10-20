package user

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
)

type Facade interface {
	RegisterUser(ctx context.Context, reqBytes []byte) ([]byte, error)
	GetProfile(ctx context.Context, reqBytes []byte, jwtPayload jwt.JwtPayload) ([]byte, error)
}
