package model

import (
	"../db"
	"fmt"
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

// ProblemWithUserStat a problem model
type ProblemWithUserStat struct {
	ID          int    	`json:"id"`
	CategoryID  int    	`json:"category_id"`
	Title       string 	`json:"title"`
	Difficulty  int    	`json:"difficulty"`
	Percent		float32	`json:"percent"`
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

// GetAllProblemsWithUserProgres return all problems in db
func GetAllProblemsWithUserProgres(uid int) ([]*ProblemWithUserStat, error) {
	statement :=
`
select 
	problem.id,
	problem.title,
	problem.categoryid,
	problem.difficulty,
	(case when submission.max is NULL THEN -1 ELSE submission.max END) as percent
from public.problem as problem
left join
(
	select
		public.submission.problem_id,
		max(public.submission.score * 100.0 / public.submission.max_score)
	from public.submission
	where public.submission.usr_id = $1
	group by public.submission.problem_id
) as submission
on problem.id = submission.problem_id
order by problem.id
`
	rows, err := db.DB.Query(statement, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	problems := make([]*ProblemWithUserStat, 0)
	for rows.Next() {
		problem := new(ProblemWithUserStat)

		err := rows.Scan(&problem.ID, &problem.Title, &problem.CategoryID, &problem.Difficulty, &problem.Percent)
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

func countAllProblems() int {
	var count int
	statement := "SELECT COUNT(*) FROM problem;"
	row := db.DB.QueryRow(statement)
	row.Scan(&count)
	return count
}

// NewTestcase add a new testcase with just input and we'll get output of testcase by using Judge0 Ou
func NewTestcase(id int, testcase Testcase) error {
	statement := `UPDATE public.problem SET testcase = array_append(testcase, ($1,$2)::TESTCASE) WHERE id=$3`
	_, err := db.DB.Exec(statement, testcase.Input, testcase.Output, id)
	if err != nil {
		return err
	}
	return nil
}

// SpecificTestcaseWithID return turncase from specific problem
func SpecificTestcaseWithID(id int) ([]*Testcase, error) {
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

// UpdateTestcase add a new testcase with just input and we'll get output of testcase by using Judge0 Ou
func UpdateTestcase(id int, index int, testcase Testcase) error {
	statement := `UPDATE public.problem SET testcase = (testcase[0:$1] || ($3,$4)::TESTCASE || testcase[$2:]) WHERE id=$5`
	_, err := db.DB.Exec(statement, index - 1, index + 1, testcase.Input, testcase.Output, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTestcase add a new testcase with just input and we'll get output of testcase by using Judge0 Ou
func DeleteTestcase(id int, index int) error {
	statement := `UPDATE public.problem SET testcase = (testcase[0:$1]  || testcase[$2:]) WHERE id=$3`
	_, err := db.DB.Exec(statement, index - 1, index + 1, id)
	if err != nil {
		return err
	}
	return nil
}
