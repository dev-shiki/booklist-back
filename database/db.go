package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("mysql", "sql12734740:45bu1mdIKc@tcp(52.76.27.242:3306)/sql12734740?allowCleartextPasswords=true")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Successfully connected to the database!")
	}
}
