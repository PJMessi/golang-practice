package model

import "time"

type User struct {
	Id        string
	Email     string
	Password  *string
	FirstName *string
	LastName  *string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
