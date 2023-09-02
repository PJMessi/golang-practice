package dbmysql

import (
	"fmt"

	"github.com/pjmessi/go-database-usage/internal/pkg/model"
)

func (dbMysql *DbMysql) CreateUser(user *model.User) error {
	res := dbMysql.db.Create(user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (dbMysql *DbMysql) IsUserEmailTaken(email string) (*bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT id FROM users WHERE email=\"%s\")", email)

	var taken bool

	res := dbMysql.db.Raw(query).Scan(&taken)
	if res.Error != nil {
		return nil, res.Error
	}

	return &taken, nil
}
