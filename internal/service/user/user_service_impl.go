package user

import (
	"fmt"

	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/internal/pkg/database"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/hash"
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
}

func NewService(db database.Db, hashUtil hash.Util, passwordUtil password.Util, uuidUtil uuid.Util) Service {
	return &ServiceImpl{
		db:           db,
		hashUtil:     hashUtil,
		passwordUtil: passwordUtil,
		uuidUtil:     uuidUtil,
	}
}

func (s *ServiceImpl) CreateUser(email string, password string) (model.User, error) {
	isEmailTaken, err := s.db.IsUserEmailTaken(email)
	if err != nil {
		return model.User{}, err
	}
	if isEmailTaken {
		return model.User{}, exception.NewAlreadyExistsFromBase(exception.Base{
			Message: fmt.Sprintf("user with the email '%s' already exists", email),
			Type:    "USER.ALREADY_EXISTS",
		})
	}

	isPwStrong, err := s.passwordUtil.IsStrong(password)
	if err != nil {
		return model.User{}, err
	}
	if !isPwStrong {
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

	err = (s.db).CreateUser(&user)
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
