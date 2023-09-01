package app

import (
	"fmt"

	db "github.com/pjmessi/go-database-usage/internal/pkg/db"
)

func StartApp() {
	fmt.Println("StartApp >>>")

	var dbInstance db.Db = db.CreateDbMysql()
	dbInstance.InitializeConnection()
	defer dbInstance.CloseConnection()

	fmt.Println("<<< StartApp")
}
