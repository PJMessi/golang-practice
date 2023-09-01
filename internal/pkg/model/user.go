package model

import "time"

type User struct {
	Id        string     `gorm:"primaryKey;column:id"`
	FirstName *string    `gorm:"column:first_name"`
	LastName  *string    `gorm:"column:last_name"`
	Email     string     `gorm:"column:email"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}
