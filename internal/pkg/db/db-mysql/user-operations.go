package dbmysql

import (
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
)

func (dbMysql *DbMysql) Create(user *model.User) {
	dbMysql.db.Create(user)
}
