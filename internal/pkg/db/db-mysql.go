package db

import (
	"fmt"
	"log"
	"time"

	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbMysql struct {
	db *gorm.DB
}

func CreateDbMysql() Db {
	return &DbMysql{}
}

func (dbMySql *DbMysql) InitializeConnection(appConfig *config.AppConfig) {
	log.Printf("PROCESSING InitializeConnection >>>")

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

	log.Printf("<<< PROCESSED InitializeConnection")
}

func (dbMySql *DbMysql) CloseConnection() {
	log.Printf("PROCESSING CloseConnection >>>")
	log.Printf("<<< PROCESSED CloseConnection")
}

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
