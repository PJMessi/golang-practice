package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/passwordutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/timeutil"
	"github.com/pjmessi/golang-practice/pkg/uuidutil"
)

type ServiceImpl struct {
	db         database.Db
	logService logger.Service
}

func NewService(
	logService logger.Service,
	db database.Db,
) Service {
	return &ServiceImpl{
		db:         db,
		logService: logService,
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

	hashedPw, err := passwordutil.Hash(password)
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
	if isPwStrong := passwordutil.IsStrong(password); !isPwStrong {
		s.logService.DebugCtx(ctx, "user did not provide strong password")

		return exception.NewInvalidReqFromBase(exception.Base{
			Details: &map[string]string{"password": "password not strong"},
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

func (s *ServiceImpl) GetProfile(ctx context.Context, userId string) (model.User, error) {
	exists, user, err := s.db.GetUserById(ctx, userId)
	if err != nil {
		return model.User{}, err
	}

	if !exists {
		return model.User{}, fmt.Errorf("user with id '%s' does not exist", userId)
	}

	return user, nil
}
