package testutil

import (
	"time"

	"github.com/jaswdr/faker"
	"github.com/pjmessi/golang-practice/internal/model"
)

var Fake faker.Faker = faker.New()

func GenMockUser(partialData *model.User) model.User {
	currentTime := time.Now()

	id := Fake.UUID().V4()
	firstName := Fake.Person().FirstName()
	lastName := Fake.Person().LastName()
	password := Fake.Internet().Password()
	createdAt := Fake.Time().TimeBetween(currentTime.AddDate(0, -1, 0), currentTime)
	updatedAt := Fake.Time().TimeBetween(createdAt, currentTime)

	if partialData != nil {
		if partialData.Id != "" {
			id = partialData.Id
		}
		if partialData.FirstName != nil && *partialData.FirstName != "" {
			id = partialData.Id
		}
		if partialData.LastName != nil && *partialData.LastName != "" {
			id = partialData.Id
		}
		if partialData.Password != nil && *partialData.Password != "" {
			id = partialData.Id
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
