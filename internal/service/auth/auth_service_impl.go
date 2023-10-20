package auth

import (
	"context"
	"fmt"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/pkg/passwordutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type ServiceImpl struct {
	Service
	db         database.Db
	jwtHandler jwt.Handler
	logService logger.Service
}

func NewService(logService logger.Service, jwtHandler jwt.Handler, db database.Db) Service {
	return &ServiceImpl{
		db:         db,
		jwtHandler: jwtHandler,
		logService: logService,
	}
}

func (s *ServiceImpl) Login(ctx context.Context, email string, password string) (model.User, string, error) {
	userExists, user, err := (s.db).GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, "", err
	}
	if !userExists {
		s.logService.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' does not exist", email))
		return model.User{}, "", exception.NewUnauthenticatedFromBase(exception.Base{
			Message: "invalid credentials",
		})
	}

	if user.Password == nil {
		s.logService.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' hasn't setup his password", email))
		return model.User{}, "", exception.NewUnauthenticatedFromBase(exception.Base{
			Type: errorcode.UserPwNotSet,
		})
	}

	if isValidHash := passwordutil.IsHashCorrect(*user.Password, password); !isValidHash {
		s.logService.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' did not provide correct password", email))
		return model.User{}, "", exception.NewUnauthenticatedFromBase(exception.Base{
			Message: "invalid credentials",
		})
	}

	jwtString, err := s.jwtHandler.Generate(jwt.JwtPayload{UserId: user.Id, UserEmail: user.Email})
	if err != nil {
		return model.User{}, "", err
	}

	return user, jwtString, nil
}

func (s *ServiceImpl) VerifyJwt(ctx context.Context, jwtStr string) (jwt.JwtPayload, error) {
	isValid, jwtPayload, err := s.jwtHandler.Verify(jwtStr)
	if err != nil {
		return jwt.JwtPayload{}, err
	}

	if !isValid {
		return jwt.JwtPayload{}, exception.Unauthenticated{}
	}

	return jwtPayload, nil
}
