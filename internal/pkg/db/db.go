package db

import model "github.com/pjmessi/go-database-usage/internal/pkg/models"

type Db interface {
	InitializeConnection()
	CloseConnection()
	CreateUser() *model.User
	GetUser() *model.User
}
