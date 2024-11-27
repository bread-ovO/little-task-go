package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", "root:040131@tcp(127.0.0.1:3306)/taskDB")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

type Book struct {
	ID          int
	UserID      int
	ISBN        string
	Title       string
	Author      string
	Publisher   string
	Description string
	CoverImage  string
}
