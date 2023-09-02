package business

import (
	"time"

	"github.com/pjmessi/go-database-usage/internal/pkg/db"
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
	"github.com/pjmessi/go-database-usage/pkg/hashing"

	"github.com/google/uuid"
)

type AccountRegistrationService struct {
	db             *db.Db
	hashingUtility *hashing.HashUtility
}

func CreateAccountRegistrationService(db *db.Db, hashingUtility *hashing.HashUtility) *AccountRegistrationService {
	return &AccountRegistrationService{
		db:             db,
		hashingUtility: hashingUtility,
	}
}

func (service *AccountRegistrationService) RegisterUser(email string, password string) (*model.User, error) {
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

	(*service.db).CreateUser(&user)

	return &user, nil
}
