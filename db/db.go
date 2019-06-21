package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//Init querying postgresql db
func Init() {
	connStr := "user=lord-tantatorn dbname=ce_boostup_db sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT username FROM grader_user WHERE id = 1")

	fmt.Println(rows)
}
