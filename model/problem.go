package model

import (
	"ce-boostup-backend/db"
)

//Testcase input and output of testcases
type Testcase struct {
	input  string
	output string
}

//Problem a problem model
type Problem struct {
	ID          int
	CategoryID  int
	Title       string
	Description string
	Difficulty  int
	CreatedAt   string //time created
	UpdatedAt   string //time updated
}

//AllProblems return all problems in db
func AllProblems() ([]*Problem, error) {
	rows, err := db.DB.Query("SELECT id,title,description,categoryID,difficulty,createdAt,updatedAt FROM problem ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	problems := make([]*Problem, 0)
	for rows.Next() {
		problem := new(Problem)

		err := rows.Scan(&problem.ID, &problem.Title, &problem.Description, &problem.CategoryID, &problem.Difficulty, &problem.CreatedAt, &problem.UpdatedAt)
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

//SpecificProblemWithID return specific problem
func SpecificProblemWithID(id int) (*Problem, error) {
	statement := `SELECT id,title,description,categoryID,difficulty,createdAt,updatedAt FROM problem WHERE id=$1`
	row := db.DB.QueryRow(statement, id)

	problem := new(Problem)

	err := row.Scan(&problem.ID, &problem.Title, &problem.Description, &problem.CategoryID, &problem.Difficulty, &problem.CreatedAt, &problem.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return problem, nil
}
