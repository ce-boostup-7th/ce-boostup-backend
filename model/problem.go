package model

import (
	"ce-boostup-backend/db"
	"fmt"
	"strings"
)

//Testcase input and output of testcases
type Testcase struct {
	Input  string
	Output string
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

//NewProblem add new problem
func NewProblem(title string, categoryID int, difficulty int) error {
	statement := `INSERT INTO problem (title,categoryID,difficulty) VALUES ($1,$2,$3)`
	_, err := db.DB.Exec(statement, title, categoryID, difficulty)
	if err != nil {
		return err
	}
	return nil
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

//SpecificTestcaseWithID return turncase from specific problem
func SpecificTestcaseWithID(id int) ([]*Testcase, error) {
	statement := `SELECT UNNEST(testcase) FROM problem WHERE id=$1`
	rows, err := db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}

	testcases := make([]*Testcase, 0)
	for rows.Next() {
		testcase := new(Testcase)
		var s, input, output string
		err := rows.Scan(&s)
		seperateString(s, &input, &output)
		if err != nil {
			return nil, err
		}
		testcase.Input = input
		testcase.Output = output
		testcases = append(testcases, testcase)
	}

	return testcases, nil
}

//UpdateProblem update problem datat
func UpdateProblem(problem Problem) error {
	statement := `UPDATE problem SET title=$1,description=$2,categoryID=$3,difficulty=$4,updatedat=CURRENT_TIMESTAMP WHERE id=$5`
	_, err := db.DB.Exec(statement, problem.Title, problem.Description, problem.CategoryID, problem.Difficulty, problem.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//DeleteAllProblems cleans all problem
func DeleteAllProblems() error {
	statement := "DELETE FROM problem; ALTER SEQUENCE problem_id_seq RESTART WITH 1;"
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

//DeleteProblemWithSpecificID delete problem by id
func DeleteProblemWithSpecificID(id int) error {
	statement := fmt.Sprintf("DELETE FROM problem WHERE id=%d ; ALTER SEQUENCE problem_id_seq RESTART WITH 1;", id)
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func seperateString(str string, str1 *string, str2 *string) {
	str = strings.Replace(str, "(", "", 1)
	str = strings.Replace(str, ")", "", 1)
	s := strings.Split(str, ",")
	*str1, *str2 = s[0], s[1]
}
