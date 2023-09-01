package dbmysql

import (
	"fmt"
)

func (dbmysql *DbMysql) IsHealthy() bool {
	sqlQuery := "SELECT COUNT(*) FROM users;"

	var result int

	if err := dbmysql.db.Raw(sqlQuery).Scan(&result).Error; err != nil {
		fmt.Printf("health check query failed with error: %v\n", err)
		return false
	}

	return true
}
