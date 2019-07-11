package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" //justify
)

//DB a pointer to sql database
var DB *sql.DB

//Init postgresql db
func Init() {
	connStr := fmt.Sprintf("user=%s password=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
	DB.SetMaxIdleConns(5)
}
