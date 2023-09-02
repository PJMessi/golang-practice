package model

import "time"

type User struct {
	Id        string     `gorm:"primaryKey;column:id"`
	Email     string     `gorm:"column:email"`
	Password  *string    `gorm:"column:password"`
	FirstName *string    `gorm:"column:first_name"`
	LastName  *string    `gorm:"column:last_name"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}
