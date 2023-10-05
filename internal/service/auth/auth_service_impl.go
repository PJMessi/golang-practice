package auth

import (
	"context"
	"fmt"

	"github.com/pjmessi/go-database-usage/internal/errorcode"
	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/internal/pkg/database"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/hash"
	"github.com/pjmessi/go-database-usage/pkg/jwt"
	"github.com/pjmessi/go-database-usage/pkg/logger"
)

type ServiceImpl struct {
	Service
	db         database.Db
	hashUtil   hash.Util
	jwtUtil    jwt.Util
	loggerUtil logger.Util
}

func NewService(loggerUtil logger.Util, jwtUtil jwt.Util, db database.Db, hashUtil hash.Util) Service {
	return &ServiceImpl{
		db:         db,
		hashUtil:   hashUtil,
		jwtUtil:    jwtUtil,
		loggerUtil: loggerUtil,
	}
}

func (s *ServiceImpl) Login(ctx context.Context, email string, password string) (model.User, string, error) {
	userExists, user, err := (s.db).GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, "", err
	}
	if !userExists {
		s.loggerUtil.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' does not exist", email))
		return model.User{}, "", exception.NewUnauthenticated()
	}

	if user.Password == nil {
		s.loggerUtil.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' hasn't setup his password", email))
		return model.User{}, "", exception.NewUnauthenticatedFromBase(exception.Base{
			Type: errorcode.UserPwNotSet,
		})
	}

	if isValidHash := s.hashUtil.VerifyHash(*user.Password, password); !isValidHash {
		s.loggerUtil.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' did not provide correct password", email))
		return model.User{}, "", exception.NewUnauthenticated()
	}

	jwtString, err := s.jwtUtil.Generate(user.Id, user.Email)
	if err != nil {
		return model.User{}, "", err
	}

	return user, jwtString, nil
}
