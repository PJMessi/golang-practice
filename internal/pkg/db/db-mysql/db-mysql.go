package dbmysql

import (
	"fmt"
	"log"

	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbMysql struct {
	db *gorm.DB
}

func CreateDbMysql() db.Db {
	return &DbMysql{}
}

func (dbMySql *DbMysql) InitializeConnection(appConfig *config.AppConfig) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			appConfig.DB_USER,
			appConfig.DB_PASSWORD,
			appConfig.DB_HOST,
			appConfig.DB_PORT,
			appConfig.DB_DATABASE,
		),
		DefaultStringSize:         255,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("error while initializing mysql connection: %v", err)
	}

	dbMySql.db = db
}

func (dbMySql *DbMysql) CloseConnection() {
}
