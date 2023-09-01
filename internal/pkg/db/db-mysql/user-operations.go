package dbmysql

import (
	"time"

	"github.com/pjmessi/go-database-usage/internal/pkg/model"
)

func (dbMySql *DbMysql) CreateUser() *model.User {
	user := &model.User{
		Id:        "one",
		Email:     "pjmessi25@icloud.com",
		CreatedAt: time.Now(),
	}
	dbMySql.db.Create(user)
	return user
}

func (dbMySql *DbMysql) GetUser() *model.User {
	var user model.User
	dbMySql.db.First(&user)
	return &user
}
