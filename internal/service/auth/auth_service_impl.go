package auth

import (
	"context"
	"fmt"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type ServiceImpl struct {
	Service
	db           database.Db
	jwtHandler   jwt.Handler
	loggerUtil   logger.Util
	passwordUtil password.Util
}

func NewService(loggerUtil logger.Util, jwtHandler jwt.Handler, db database.Db, passwordUtil password.Util) Service {
	return &ServiceImpl{
		db:           db,
		jwtHandler:   jwtHandler,
		loggerUtil:   loggerUtil,
		passwordUtil: passwordUtil,
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

	if isValidHash := s.passwordUtil.IsHashCorrect(*user.Password, password); !isValidHash {
		s.loggerUtil.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' did not provide correct password", email))
		return model.User{}, "", exception.NewUnauthenticated()
	}

	jwtString, err := s.jwtHandler.Generate(user.Id, user.Email)
	if err != nil {
		return model.User{}, "", err
	}

	return user, jwtString, nil
}
