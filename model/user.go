package model

import (
	"ce-boostup-backend/db"
	"fmt"
)

//User a grader user model
type User struct {
	ID       int
	Username string
	Password string
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

		err := rows.Scan(&user.ID, &user.Username, &user.Password)
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

	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//UpdateUser update user data
func UpdateUser(usr User) error {
	statement := `UPDATE grader_user SET username=$1, password=$2 WHERE id=$3;`
	_, err := db.DB.Exec(statement, usr.Username, usr.Password, usr.ID)
	if err != nil {
		fmt.Println("hola")
		return err
	}
	return nil
}

//DeleteAllUsers clean users
func DeleteAllUsers() error {
	statement := "DELETE FROM grader_user; ALTER SEQUENCE grader_user_id_seq RESTART WITH 1;"
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

//DeleteUserWithSpecificID delete user by id
func DeleteUserWithSpecificID(id int) error {
	statement := fmt.Sprintf("DELETE FROM grader_user WHERE id=%d ; ALTER SEQUENCE grader_user_id_seq RESTART WITH 1;", id)
	_, err := db.DB.Exec(statement)

	if err != nil {
		return err
	}
	return nil
}
