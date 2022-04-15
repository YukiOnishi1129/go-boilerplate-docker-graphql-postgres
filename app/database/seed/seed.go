package main

import (
	"database/sql"
	"fmt"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/seed/seeders"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
		panic(err)
	}

	conn, dbErr := connectDB()
	if dbErr != nil {
		panic(dbErr)
	}
	seedErr := seeders.CreateTestData(conn)
	if seedErr != nil {
		panic(seedErr)
	}
}

func connectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv(("MYSQL_USER")), os.Getenv(("MYSQL_PASSWORD")), os.Getenv(("MYSQL_PORT")), os.Getenv(("MYSQL_DATABASE")))
	db, dbErr := sql.Open("mysql", dsn)
	if dbErr != nil {
		return nil, dbErr
	}
	return db, nil
}
