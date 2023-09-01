package db

import (
	"log"
	"time"

	model "github.com/pjmessi/go-database-usage/internal/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbMysql struct {
	db *gorm.DB
}

func CreateDbMysql() Db {
	return &DbMysql{}
}

func (dbMySql *DbMysql) InitializeConnection() {
	log.Printf("PROCESSING InitializeConnection >>>")

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:miferia@tcp(127.0.0.1:3307)/go-database-test?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                       // default size for string fields
		DisableDatetimePrecision:  true,                                                                                      // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                      // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                      // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                     // auto configure based on currently MySQL version
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
