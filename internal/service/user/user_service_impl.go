package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/timeutil"
	"github.com/pjmessi/golang-practice/pkg/uuidutil"
)

type ServiceImpl struct {
	db           database.Db
	passwordUtil password.Util
	logService   logger.Service
}

func NewService(
	logService logger.Service,
	db database.Db,
	passwordUtil password.Util,
) Service {
	return &ServiceImpl{
		db:           db,
		passwordUtil: passwordUtil,
		logService:   logService,
	}
}

func (s *ServiceImpl) CreateUser(ctx context.Context, email string, password string) (model.User, error) {
	lowercaseEmail := strings.ToLower(email)

	if err := s.ensureStrongPw(ctx, password); err != nil {
		return model.User{}, err
	}

	if err := s.ensureEmailNotUsed(ctx, lowercaseEmail); err != nil {
		return model.User{}, err
	}

	hashedPw, err := s.passwordUtil.Hash(password)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.createUser(lowercaseEmail, hashedPw)
	if err != nil {
		return model.User{}, err
	}

	err = s.db.SaveUser(ctx, &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *ServiceImpl) ensureStrongPw(ctx context.Context, password string) error {
	if isPwStrong := s.passwordUtil.IsStrong(password); !isPwStrong {
		s.logService.DebugCtx(ctx, "user did not provide strong password")

		return exception.NewInvalidReqFromBase(exception.Base{
			Message: "password is not strong enough",
		})
	}

	return nil
}

func (s *ServiceImpl) ensureEmailNotUsed(ctx context.Context, email string) error {
	isEmailTaken, err := s.db.IsUserEmailTaken(ctx, email)
	if err != nil {
		return err
	}

	if !isEmailTaken {
		return nil
	}

	s.logService.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' already exists", email))

	return exception.NewAlreadyExistsFromBase(exception.Base{
		Message: fmt.Sprintf("user with the email '%s' already exists", email),
		Type:    errorcode.UserAlreadyExist,
	})
}

func (s *ServiceImpl) createUser(email string, hashedPw string) (model.User, error) {
	uuidStr, err := uuidutil.GenUuidV4()
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Id:        uuidStr,
		Email:     email,
		Password:  &hashedPw,
		CreatedAt: timeutil.GetCurrentTime(),
	}, nil
}
