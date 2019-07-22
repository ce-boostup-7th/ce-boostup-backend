package model

import (
	"../conversion"
	"../db"
	"../judge0"
	"database/sql"
)

// OuSubmission a model for submission
type OuSubmission struct {
	SubmissionID 	int     `json:"submission_id" form:"submission_id"`
	UserID       	int     `json:"user_id" form:"user_id"`
	ProblemID    	int     `json:"problem_id" form:"problem_id"`
	LanguageID   	int     `json:"language_id" form:"language_id"`
	Src          	string  `json:"src" form:"src"`
	SubmittedAt  	string  `json:"submitted_at" form:"submitted_at"`
	Score        	int     `json:"score" form:"score"`
	MaxScore     	int     `json:"max_score" form:"max_score"`
	Runtime      	float64 `json:"runtime" form:"runtime"`
	MemoryUsage  	int     `json:"memory_usage" form:"memory_usage"`
	Results			string 	`json:"results" form:"results"`
	CompileOutput 	string	`json:"compile_output" form:"compile_output"`
}

// OuNewSubmission create a new submission
func OuNewSubmission(userID int, problemID int, languageID int, src string) (int, error) {
	

	testcases, err := OuSpecificTestcaseWithID(problemID)
	if err != nil {
		return -1, err
	}

	score := 0
	runtime := 0.0
	memory := 0
	resultsArr := make([]byte, len(testcases))
	CompileOutput := ""

	ch := make(chan *judge0.OuResult)
	chIndex := make(chan int)

	for i := range testcases {
		go func(testcase *Testcase, i int) {
			ch <- judge0.OuSubmit(languageID, src, testcase.Input, testcase.Output)
			chIndex <- i
		}(testcases[i], i)
	}

	for range testcases {
		result := <- ch
		index := <- chIndex
		memory += result.Memory
		runtime += conversion.StringToFloat(result.Time)
		switch result.Status.ID {
		case 3: // Accepted
			score++
			resultsArr[index] = 'P'
		case 4: // Wrong Answer
			resultsArr[index] = '-'
		case 5: // Time Limit Exceeded
			resultsArr[index] = 'T'
		case 6: // Compilation Error
			resultsArr[index] = 'C'
		case 13: // Internal Error
			resultsArr[index] = 'I'
		default:
			resultsArr[index] = 'X'
		}
		CompileOutput = result.CompileOutput
	}

	length := len(testcases)
	runtime = runtime / float64(length)
	memory = memory / length

	results := string(resultsArr)

	statement := `INSERT INTO public.submission
	(usr_id, problem_id, lang_id, src, score, runtime, memory_usage, max_score, results, compile_output)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING submission_id`
	row := db.DB.QueryRow(statement, userID, problemID, languageID, src, score, runtime, memory, length, results , CompileOutput)

	lastInsertID := -1
	err = row.Scan(&lastInsertID)
	if err != nil {
		return -1, err
	}

	err = collectScore(userID)
	if err != nil {
		return -1, err
	}

	return lastInsertID, nil
}

var baseSQL = `SELECT submission_id, src, usr_id, problem_id, lang_id, submittedat, score, max_score, runtime, memory_usage, results, compile_output
	FROM public.submission`

// OuAllSubmissions get all submissions
func OuAllSubmissions() ([]*OuSubmission, error) {
	rows, err := db.DB.Query( baseSQL + " ORDER BY submission_id")
	if err != nil {
		return nil, err
	}

	submissions, err := scansSubmissionRow(rows)
	if err != nil {
		return nil, err
	}

	return submissions, nil
}

// OuAllSubmissionsFilteredByUserID filter all user that filtered by userID
func OuAllSubmissionsFilteredByUserID(uid int) ([]*OuSubmission, error) {
	rows, err := db.DB.Query(baseSQL + " WHERE usr_id=$1 ORDER BY submission_id", uid)
	if err != nil {
		return nil, err
	}

	submissions, err := scansSubmissionRow(rows)
	if err != nil {
		return nil, err
	}

	return submissions, nil
}

// OuSpecificSubmission return a specific submission by id
func OuSpecificSubmission(id int) (*OuSubmission, error) {
	row := db.DB.QueryRow(baseSQL + " WHERE submission_id=$1", id)

	submission, err := scanSubmissionRow(row)
	if err != nil {
		return nil, err
	}
	
	return submission, nil
}

func scanSubmissionRow(row *sql.Row) (*OuSubmission, error) {
	submission := new(OuSubmission)

	var results sql.NullString
	var compileOutput sql.NullString
	
	err := row.Scan(&submission.SubmissionID, 
		&submission.Src, 
		&submission.UserID, 
		&submission.ProblemID, 
		&submission.LanguageID, 
		&submission.SubmittedAt, 
		&submission.Score, 
		&submission.MaxScore, 
		&submission.Runtime, 
		&submission.MemoryUsage,
		&results,
		&compileOutput,
	)

	if results.Valid {
		submission.Results = results.String
	}
	if compileOutput.Valid {
		submission.CompileOutput = compileOutput.String
	}

	if err != nil {
		return nil, err
	}
	return submission, nil
}

func scansSubmissionRow(rows *sql.Rows) ([]*OuSubmission, error) {
	submissions := make([]*OuSubmission, 0)

	defer rows.Close()

	for rows.Next() {
		submission := new(OuSubmission)

		var results sql.NullString
		var compileOutput sql.NullString

		err := rows.Scan(&submission.SubmissionID, 
			&submission.Src,
			&submission.UserID,
			&submission.ProblemID,
			&submission.LanguageID,
			&submission.SubmittedAt,
			&submission.Score,
			&submission.MaxScore,
			&submission.Runtime,
			&submission.MemoryUsage,
			&results,
			&compileOutput,
		)

		if results.Valid {
			submission.Results = results.String
		}
		if compileOutput.Valid {
			submission.CompileOutput = compileOutput.String
		}

		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return submissions, nil
}