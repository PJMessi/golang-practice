package auth

import (
	"context"

	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/internal/pkg/database"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/hash"
	"github.com/pjmessi/go-database-usage/pkg/jwt"
)

type ServiceImpl struct {
	Service
	db       database.Db
	hashUtil hash.Util
	jwtUtil  jwt.Util
}

func NewService(jwtUtil jwt.Util, db database.Db, hashUtil hash.Util) Service {
	return &ServiceImpl{
		db:       db,
		hashUtil: hashUtil,
		jwtUtil:  jwtUtil,
	}
}

func (s *ServiceImpl) Login(ctx context.Context, email string, password string) (model.User, string, error) {
	userExists, user, err := (s.db).GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, "", err
	}
	if !userExists {
		return model.User{}, "", exception.NewUnauthenticated()
	}

	if user.Password == nil {
		return model.User{}, "", exception.NewUnauthenticatedFromBase(exception.Base{
			Type: "UNAUTHENTICATED.PASSWORD_NOT_SET",
		})
	}

	if isValidHash := s.hashUtil.VerifyHash(*user.Password, password); !isValidHash {
		return model.User{}, "", exception.NewUnauthenticated()
	}

	jwtString, err := s.jwtUtil.Generate(user.Id, user.Email)
	if err != nil {
		return model.User{}, "", err
	}

	return user, jwtString, nil
}
