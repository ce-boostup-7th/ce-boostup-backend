package model

import (
	"ce-boostup-backend/db"
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
	rows, err := db.DB.Query("SELECT * FROM grader_user")
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
