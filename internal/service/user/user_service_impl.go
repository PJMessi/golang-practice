package user

import (
	"context"
	"fmt"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/timeutil"
	"github.com/pjmessi/golang-practice/pkg/uuid"
)

type ServiceImpl struct {
	Service
	db           database.Db
	passwordUtil password.Util
	uuidUtil     uuid.Util
	loggerUtil   logger.Util
}

func NewService(
	loggerUtil logger.Util,
	db database.Db,
	passwordUtil password.Util,
	uuidUtil uuid.Util,
) Service {
	return &ServiceImpl{
		db:           db,
		passwordUtil: passwordUtil,
		uuidUtil:     uuidUtil,
		loggerUtil:   loggerUtil,
	}
}

func (s *ServiceImpl) CreateUser(ctx context.Context, email string, password string) (model.User, error) {
	isEmailTaken, err := s.db.IsUserEmailTaken(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	if isEmailTaken {
		s.loggerUtil.DebugCtx(ctx, fmt.Sprintf("user with the email '%s' already exists", email))
		return model.User{}, exception.NewAlreadyExistsFromBase(exception.Base{
			Message: fmt.Sprintf("user with the email '%s' already exists", email),
			Type:    errorcode.UserAlreadyExist,
		})
	}

	if isPwStrong := s.passwordUtil.IsStrong(password); !isPwStrong {
		s.loggerUtil.DebugCtx(ctx, "user did not provide strong password")
		return model.User{}, exception.NewInvalidReqFromBase(exception.Base{
			Message: "password is not strong enough",
		})
	}

	hashedPw, err := s.passwordUtil.Hash(password)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.createUserModel(email, hashedPw)
	if err != nil {
		return model.User{}, err
	}

	err = (s.db).CreateUser(ctx, &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *ServiceImpl) createUserModel(email string, hashedPw string) (model.User, error) {
	uuidStr, err := s.uuidUtil.GenUuidV4()
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Id:        uuidStr,
		Email:     email,
		Password:  &hashedPw,
		CreatedAt: timeutil.GetCurrentDateTimeStr(),
	}, nil
}
