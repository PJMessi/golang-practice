package model

import "time"

type User struct {
	Id        string
	FirstName *string
	LastName  *string
	Email     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
