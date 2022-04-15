package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

//Init DB接続処理
func Init() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv(("MYSQL_USER")), os.Getenv(("MYSQL_PASSWORD")), os.Getenv(("MYSQL_HOST")), os.Getenv(("MYSQL_PORT")), os.Getenv(("MYSQL_DATABASE")))
	db, dbErr := sql.Open("mysql", dsn)
	if dbErr != nil {
		return nil, dbErr
	}
	return db, nil
}
