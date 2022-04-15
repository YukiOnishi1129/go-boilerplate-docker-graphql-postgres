package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

//Init DB接続処理
func Init() (*sql.DB, error) {
	dsn := fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, dbErr := sql.Open("postgres", dsn)
	if dbErr != nil {
		return nil, dbErr
	}
	return db, nil
}
