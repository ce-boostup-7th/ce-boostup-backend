package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" //justify
)

//DB a pointer to sql database
var DB *sql.DB

//Init postgresql db
func Init() {
	connStr := "user=lord-tantatorn password=pass port=5432 dbname=ce_boostup_db sslmode=disable"
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
