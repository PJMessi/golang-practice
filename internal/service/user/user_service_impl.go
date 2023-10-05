package user

import (
	"context"
	"fmt"

	"github.com/pjmessi/go-database-usage/internal/errorcode"
	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/internal/pkg/database"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/hash"
	"github.com/pjmessi/go-database-usage/pkg/logger"
	"github.com/pjmessi/go-database-usage/pkg/password"
	"github.com/pjmessi/go-database-usage/pkg/timeutil"
	"github.com/pjmessi/go-database-usage/pkg/uuid"
)

type ServiceImpl struct {
	Service
	db           database.Db
	hashUtil     hash.Util
	passwordUtil password.Util
	uuidUtil     uuid.Util
	loggerUtil   logger.Util
}

func NewService(
	loggerUtil logger.Util,
	db database.Db,
	hashUtil hash.Util,
	passwordUtil password.Util,
	uuidUtil uuid.Util,
) Service {
	return &ServiceImpl{
		db:           db,
		hashUtil:     hashUtil,
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

	isPwStrong, err := s.passwordUtil.IsStrong(password)
	if err != nil {
		return model.User{}, err
	}
	if !isPwStrong {
		s.loggerUtil.DebugCtx(ctx, "user did not provide strong password")
		return model.User{}, exception.NewInvalidReqFromBase(exception.Base{
			Message: "password is not strong enough",
		})
	}

	hashedPw, err := s.hashUtil.GenerateHash(password)
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
