package model

import (
	"../db"
	"fmt"
	"strings"
)

//Testcase input and output of testcases
type Testcase struct {
	Input  string `json:"input" form:"input"`
	Output string `json:"output" form:"output"`
}

//Problem a problem model
type Problem struct {
	ID          int    `json:"id" form:"id"`
	CategoryID  int    `json:"category_id" form:"category_id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Difficulty  int    `json:"difficulty" form:"difficulty"`
	CreatedAt   string `json:"created_at" form:"created_at"` //time created
	UpdatedAt   string `json:"updated_at" form:"updated_at"` //time updated
}

//NewProblem add new problem
func NewProblem(title string, categoryID int, difficulty int, description string) (*int, error) {
	var problemID int

	statement := `INSERT INTO problem (title,categoryID,difficulty,description) VALUES ($1,$2,$3,$4) RETURNING id`
	row := db.DB.QueryRow(statement, title, categoryID, difficulty, description)
	err := row.Scan(&problemID)
	if err != nil {
		return nil, err
	}
	return &problemID, nil
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

// NewTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func NewTestcase(id int, testcase Testcase) error {
	statement := `UPDATE problem SET testcase=array_append(testcase,($1,$2)::TESTCASE) WHERE id=$3`
	_, err := db.DB.Exec(statement, testcase.Input+" ", testcase.Output+" ", id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateProblem update problem datat
func UpdateProblem(problem Problem) error {
	statement := `UPDATE problem SET title=$1,description=$2,categoryID=$3,difficulty=$4,updatedat=CURRENT_TIMESTAMP WHERE id=$5`
	_, err := db.DB.Exec(statement, problem.Title, problem.Description, problem.CategoryID, problem.Difficulty, problem.ID)
	if err != nil {
		return err
	}
	return nil
}

//DeleteAllProblems cleans all problem
func DeleteAllProblems() error {
	statement := "DELETE FROM problem;"
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

//DeleteProblemWithSpecificID delete problem by id
func DeleteProblemWithSpecificID(id int) error {
	statement := fmt.Sprintf("DELETE FROM problem WHERE id=%d ;", id)
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

func countAllProblems() int {
	var count int
	statement := "SELECT COUNT(*) FROM problem;"
	row := db.DB.QueryRow(statement)
	row.Scan(&count)
	return count
}

// ---------- OU version ----------

// OuNewTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func OuNewTestcase(id int, testcase Testcase) error {
	statement := `UPDATE problem SET testcase=array_append(testcase,($1,$2)::TESTCASE) WHERE id=$3`
	_, err := db.DB.Exec(statement, testcase.Input, testcase.Output, id)
	if err != nil {
		return err
	}
	return nil
}

// OuSpecificTestcaseWithID return turncase from specific problem
func OuSpecificTestcaseWithID(id int) ([]*Testcase, error) {
	statement := `select (m.u).input, (m.u).output from (SELECT UNNEST(testcase) as u FROM public.problem WHERE id=$1) as m`
	rows, err := db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}

	testcases := make([]*Testcase, 0)
	for rows.Next() {
		testcase := new(Testcase)
		err := rows.Scan(&testcase.Input, &testcase.Output)
		if err != nil {
			return nil, err
		}
		testcases = append(testcases, testcase)
	}

	return testcases, nil
}

// OuUpdateTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func OuUpdateTestcase(id int, index int, testcase Testcase) error {
	statement := `UPDATE public.problem SET testcase = (testcase[0:$1] || ($3,$4)::TESTCASE || testcase[$2:]) WHERE id=$5`
	_, err := db.DB.Exec(statement, index - 1, index + 1, testcase.Input, testcase.Output, id)
	if err != nil {
		return err
	}
	return nil
}

// OuDeleteTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func OuDeleteTestcase(id int, index int) error {
	statement := `UPDATE public.problem SET testcase = (testcase[0:$1]  || testcase[$2:]) WHERE id=$3`
	_, err := db.DB.Exec(statement, index - 1, index + 1, id)
	if err != nil {
		return err
	}
	return nil
}