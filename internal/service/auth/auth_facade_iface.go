package auth

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
)

type Facade interface {
	Login(ctx context.Context, reqBytes []byte) ([]byte, error)
	VerifyJwt(ctx context.Context, jwtStr string) (jwt.JwtPayload, error)
}
