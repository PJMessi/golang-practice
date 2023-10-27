package testutil

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	createdAt := Fake.Time().TimeBetween(currentTime.AddDate(0, -1, 0), currentTime).UTC()
	updatedAt := Fake.Time().TimeBetween(createdAt, currentTime).UTC()

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
		if (partialData.CreatedAt != time.Time{}) {
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

func GenMockLoginApiReq(partialData *model.LoginApiReq) model.LoginApiReq {
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

func GenMockRegUserApiReq(partialData *model.UserRegApiReq) model.UserRegApiReq {
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

func SetupTestUser(testServerUrl string) model.LoginApiRes {
	// registering new user
	url := fmt.Sprintf("%s/users/registration", testServerUrl)
	email := Fake.Internet().Email()
	password := "Password123!"
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// logging in the new user
	url = fmt.Sprintf("%s/auth/login", testServerUrl)
	reqBody = []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	responseBodyByte, _ := io.ReadAll(resp.Body)
	responseBody := model.LoginApiRes{}
	_ = json.Unmarshal(responseBodyByte, &responseBody)

	return responseBody
}

func GetMockAppConfig(appConf *config.AppConfig) config.AppConfig {
	finalAppConfig := config.AppConfig{
		APP_PORT:                     "3000",
		DB_HOST:                      "localhost",
		DB_PORT:                      "3006",
		DB_DATABASE:                  "go_test",
		DB_USER:                      Fake.Internet().User(),
		DB_PASSWORD:                  Fake.Internet().Password(),
		JWT_SECRET:                   Fake.RandomStringWithLength(10),
		JWT_EXPIRATION_TIME:          "1d",
		NATS_URL:                     "nats://127.0.0.1:4222",
		NATS_STREAM:                  "GO_STREAM",
		NATS_EVENT_USER_REGISTRATION: "EVENT.USER.NEW",
	}

	if appConf != nil {
		if appConf.APP_PORT != "" {
			finalAppConfig.APP_PORT = appConf.APP_PORT
		}
		if appConf.DB_HOST != "" {
			finalAppConfig.DB_HOST = appConf.DB_HOST
		}
		if appConf.DB_PORT != "" {
			finalAppConfig.DB_PORT = appConf.DB_PORT
		}
		if appConf.DB_DATABASE != "" {
			finalAppConfig.DB_DATABASE = appConf.DB_DATABASE
		}
		if appConf.DB_USER != "" {
			finalAppConfig.DB_USER = appConf.DB_USER
		}
		if appConf.DB_PASSWORD != "" {
			finalAppConfig.DB_PASSWORD = appConf.DB_PASSWORD
		}
		if appConf.JWT_SECRET != "" {
			finalAppConfig.JWT_SECRET = appConf.JWT_SECRET
		}
		if appConf.JWT_EXPIRATION_TIME != "" {
			finalAppConfig.JWT_EXPIRATION_TIME = appConf.JWT_EXPIRATION_TIME
		}
		if appConf.NATS_URL != "" {
			finalAppConfig.NATS_URL = appConf.NATS_URL
		}
		if appConf.NATS_STREAM != "" {
			finalAppConfig.NATS_STREAM = appConf.NATS_STREAM
		}
		if appConf.NATS_EVENT_USER_REGISTRATION != "" {
			finalAppConfig.NATS_EVENT_USER_REGISTRATION = appConf.NATS_EVENT_USER_REGISTRATION
		}
	}

	return finalAppConfig
}
