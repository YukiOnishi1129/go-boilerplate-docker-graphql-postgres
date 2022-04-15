package main

import (
	"database/sql"
	"fmt"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/database/seed/seeders"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	dsn := fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_PORT"), "localhost", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, dbErr := sql.Open("postgres", dsn)
	if dbErr != nil {
		return nil, dbErr
	}
	return db, nil
}
