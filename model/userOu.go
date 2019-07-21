package model

import (
	"../db"
)

// ---------- OU version ----------

//IDPasswordByUsername get password by username
func IDPasswordByUsername(username string) (*int, *string, error) {
	statement := `SELECT id, password FROM public.grader_user WHERE username=$1`
	row := db.DB.QueryRow(statement, username)

	var id int
	var password string
	err := row.Scan(&id, &password)
	if err != nil {
		return nil, nil, err
	}
	return &id, &password, nil
}
