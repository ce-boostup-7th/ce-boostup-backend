package model

import (
	"../db"
	"fmt"
)

//User a grader user model
type User struct {
	ID       int     `json:"id" form:"id"`
	Username string  `json:"username" form:"username"`
	Password string  `json:"password" form:"password"`
	Score    float64 `json:"score" form:"score"`
}

//NewUser add a new user
func NewUser(username string, password string) error {
	statement := `INSERT INTO grader_user (username,password) VALUES ($1,$2)`
	_, err := db.DB.Exec(statement, username, password)
	if err != nil {
		return err
	}
	return nil
}

//AllUsers Get all users info from db
func AllUsers() ([]*User, error) {
	rows, err := db.DB.Query("SELECT * FROM grader_user ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)

		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Score)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

//SpecificUserWithID return specific user from db by using id
func SpecificUserWithID(id int) (*User, error) {
	statement := `SELECT * FROM grader_user WHERE id=$1`
	row := db.DB.QueryRow(statement, id)

	user := new(User)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Score)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//UpdateUser update user credential data
func UpdateUser(usr User) error {
	statement := `UPDATE grader_user SET username=$1, password=$2 WHERE id=$3;`
	_, err := db.DB.Exec(statement, usr.Username, usr.Password, usr.ID)
	if err != nil {
		return err
	}
	return nil
}

//DeleteAllUsers clean user
func DeleteAllUsers() error {
	statement := "DELETE FROM grader_user;"
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

//DeleteUserWithSpecificID delete user by id
func DeleteUserWithSpecificID(id int) error {
	statement := fmt.Sprintf("DELETE FROM grader_user WHERE id=%d ;", id)
	_, err := db.DB.Exec(statement)

	if err != nil {
		return err
	}
	return nil
}

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

//IsUserExist check if user exists
func IsUserExist(username string) (*bool, error) {
	var isExist bool
	var str string

	statement := `SELECT EXISTS(SELECT * FROM grader_user WHERE username=$1)`
	row := db.DB.QueryRow(statement, username)
	err := row.Scan(&str)

	if err != nil {
		return nil, err
	}

	if str == "true" {
		isExist = true
	} else if str == "false" {
		isExist = false
	}
	return &isExist, nil
}

//IDByUsername get user id by username
func IDByUsername(username string) (*int, error) {
	statement := `SELECT id FROM grader_user WHERE username=$1`
	row := db.DB.QueryRow(statement, username)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
