package testutil

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jaswdr/faker"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/internal/model"
)

var Fake faker.Faker = faker.New()

func GenMockUser(partialData *model.User) model.User {
	currentTime := time.Now()

	id := Fake.UUID().V4()
	email := strings.ToLower(Fake.Internet().Email())
	firstName := Fake.Person().FirstName()
	lastName := Fake.Person().LastName()
	password := Fake.Internet().Password()
	createdAt := Fake.Time().TimeBetween(currentTime.AddDate(0, -1, 0), currentTime)
	updatedAt := Fake.Time().TimeBetween(createdAt, currentTime)

	if partialData != nil {
		if partialData.Id != "" {
			id = partialData.Id
		}
		if partialData.Email != "" {
			email = partialData.Email
		}
		if partialData.FirstName != nil && *partialData.FirstName != "" {
			firstName = *partialData.FirstName
		}
		if partialData.LastName != nil && *partialData.LastName != "" {
			lastName = *partialData.LastName
		}
		if partialData.Password != nil && *partialData.Password != "" {
			password = *partialData.Password
		}
		if (partialData.CreatedAt == time.Time{}) {
			createdAt = partialData.CreatedAt
		}
		if partialData.UpdatedAt != nil && (*partialData.UpdatedAt != time.Time{}) {
			updatedAt = *partialData.UpdatedAt
		}
	}

	return model.User{
		Id:        id,
		Email:     email,
		FirstName: &firstName,
		Password:  &password,
		LastName:  &lastName,
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}
}

func GenLoginApiReq(partialData *model.LoginApiReq) model.LoginApiReq {
	email := Fake.Internet().Email()
	password := Fake.Internet().Password()

	if partialData != nil {
		if email != "" {
			email = partialData.Email
		}

		if password != "" {
			password = partialData.Password
		}
	}

	return model.LoginApiReq{
		Email:    email,
		Password: password,
	}
}

func GenRegUserApiReq(partialData *model.UserRegApiReq) model.UserRegApiReq {
	email := Fake.Internet().Email()
	password := Fake.Internet().Password()

	if partialData != nil {
		if email != "" {
			email = partialData.Email
		}

		if password != "" {
			password = partialData.Password
		}
	}

	return model.UserRegApiReq{
		Email:    email,
		Password: password,
	}
}

func GetTestDbCon(appConf *config.AppConfig) (*sql.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", appConf.DB_USER, appConf.DB_PASSWORD, appConf.DB_HOST, appConf.DB_PORT, appConf.DB_DATABASE)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, fmt.Errorf("testutil.GetTestDbCon(): %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
