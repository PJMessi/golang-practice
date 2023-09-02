package business

import (
	"fmt"
	"time"

	"github.com/pjmessi/go-database-usage/internal/exceptions"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
	"github.com/pjmessi/go-database-usage/pkg/hashing"
	"github.com/pjmessi/go-database-usage/pkg/password"

	"github.com/google/uuid"
)

type AccountRegistrationService struct {
	db              *db.Db
	hashingUtility  *hashing.HashUtility
	passwordUtility *password.PasswordUtility
}

func CreateAccountRegistrationService(
	db *db.Db,
	hashingUtility *hashing.HashUtility,
	passwordUtility *password.PasswordUtility,
) *AccountRegistrationService {
	return &AccountRegistrationService{
		db:              db,
		hashingUtility:  hashingUtility,
		passwordUtility: passwordUtility,
	}
}

func (service *AccountRegistrationService) RegisterUser(email string, password string) (*model.User, error) {
	isEmailTaken, err := (*service.db).IsUserEmailTaken(email)
	if err != nil {
		return nil, err
	}

	if *isEmailTaken {
		return nil, &exceptions.DuplicateException{
			Type:    "USER.DUPLICATE",
			Message: fmt.Sprintf("user with the email '%s' already exists", email),
		}
	}

	if !service.passwordUtility.IsStrong(password) {
		return nil, &exceptions.InvalidRequestException{
			Type:    "REQUEST_DATA.INVALID",
			Message: "password is not strong enough",
			Details: nil,
		}
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	hashedPw, err := service.hashingUtility.HashString(password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Id:        uuid.String(),
		Email:     email,
		Password:  hashedPw,
		CreatedAt: time.Now(),
	}

	err = (*service.db).CreateUser(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
