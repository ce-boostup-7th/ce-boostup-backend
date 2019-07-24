package model

import (
	"../conversion"
	"../db"
	"../judge0"
	"database/sql"
)

// Submission a model for submission Ou
type Submission struct {
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

// SubmissionNoSrc a model for submission Ou
type SubmissionNoSrc struct {
	SubmissionID 	int     `json:"submission_id"`
	UserID       	int     `json:"user_id"`
	ProblemID    	int     `json:"problem_id"`
	LanguageID   	int     `json:"language_id"`
	SubmittedAt  	string  `json:"submitted_at"`
	Score        	int     `json:"score"`
	MaxScore     	int     `json:"max_score"`
	Runtime      	float64 `json:"runtime"`
	MemoryUsage  	int     `json:"memory_usage"`
	Results			string 	`json:"results"`
	CompileOutput 	string	`json:"compile_output"`
}

// NewSubmission create a new submission ou
func NewSubmission(userID int, problemID int, languageID int, src string) (*Submission, error) {
	testcases, err := SpecificTestcaseWithID(problemID)
	if err != nil {
		return nil, err
	}

	submission := new(Submission)

	submission.UserID = userID
	submission.ProblemID = problemID
	submission.LanguageID = languageID
	submission.Src = src
	submission.MaxScore = len(testcases)
	
	resultsArr := make([]byte, len(testcases))

	ch := make(chan *judge0.Result)
	chIndex := make(chan int)

	for i := range testcases {
		go func(testcase *Testcase, i int) {
			ch <- judge0.Submit(languageID, src, testcase.Input, testcase.Output)
			chIndex <- i
		}(testcases[i], i)
	}

	for range testcases {
		result := <- ch
		index := <- chIndex

		if result == nil {
			resultsArr[index] = 'B'
			continue
		}

		submission.MemoryUsage += result.Memory
		submission.Runtime += conversion.StringToFloat(result.Time)
		switch result.Status.ID {
		case 3: // Accepted
			submission.Score++
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
			switch result.Status.ID {
			case 1:
				result.CompileOutput = "In Queue"
			case 2:
				result.CompileOutput = "Processing"
			case 7:
				result.CompileOutput = "Runtime Error (SIGSEGV)"
			case 8:
				result.CompileOutput = "Runtime Error (SIGXFSZ)"
			case 9:
				result.CompileOutput = "Runtime Error (SIGFPE)"
			case 10:
				result.CompileOutput = "Runtime Error (SIGABRT)"
			case 11:
				result.CompileOutput = "Runtime Error (NZEC)"
			case 12:
				result.CompileOutput = "Runtime Error (Other)"
			}
		}
		submission.CompileOutput = result.CompileOutput
	}

	submission.Runtime /= float64(submission.MaxScore)
	submission.MemoryUsage /= submission.MaxScore
	submission.Results = string(resultsArr)

	statement := `INSERT INTO public.submission
	(usr_id, problem_id, lang_id, src, score, runtime, memory_usage, max_score, results, compile_output)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING submission_id, submittedat`
	row := db.DB.QueryRow(statement,
		submission.UserID,
		submission.ProblemID,
		submission.LanguageID,
		submission.Src,
		submission.Score,
		submission.Runtime,
		submission.MemoryUsage,
		submission.MaxScore,
		submission.Results,
		submission.CompileOutput,
	)

	submission.SubmissionID = -1
	err = row.Scan(&submission.SubmissionID, &submission.SubmittedAt)
	if err != nil {
		return nil, err
	}

	err = collectScore(userID)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

var baseSQL = `SELECT submission_id, src, usr_id, problem_id, lang_id, submittedat, score, max_score, runtime, memory_usage, results, compile_output
	FROM public.submission`

// AllSubmissions get all submissions Ou
func AllSubmissions() ([]*SubmissionNoSrc, error) {
	rows, err := db.DB.Query( baseSQL + " ORDER BY submission_id")
	if err != nil {
		return nil, err
	}

	submissions, err := scansSubmissionRow(rows)
	if err != nil {
		return nil, err
	}

	submissionsNoSrc := make([]*SubmissionNoSrc, 0)

	for _, val := range submissions {
		submissionNoSrc := new(SubmissionNoSrc)
		submissionNoSrc.SubmissionID = val.SubmissionID
		submissionNoSrc.UserID = val.UserID
		submissionNoSrc.ProblemID = val.ProblemID
		submissionNoSrc.LanguageID = val.LanguageID
		submissionNoSrc.SubmittedAt = val.SubmittedAt
		submissionNoSrc.Score = val.Score
		submissionNoSrc.MaxScore = val.MaxScore
		submissionNoSrc.Runtime = val.Runtime
		submissionNoSrc.MemoryUsage = val.MemoryUsage
		submissionNoSrc.Results = val.Results
		submissionNoSrc.CompileOutput = val.CompileOutput

		submissionsNoSrc = append(submissionsNoSrc, submissionNoSrc)
	}

	return submissionsNoSrc, nil
}

// AllSubmissionsFilteredByUserID filter all user that filtered by userID Ou
func AllSubmissionsFilteredByUserID(uid int) ([]*Submission, error) {
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

// LastUserSubmissionsFilteredByProblemID filter all user that filtered by userID Ou
func LastUserSubmissionsFilteredByProblemID(uid int, pid int) (*Submission, error) {
	row := db.DB.QueryRow(baseSQL + ` 
where public.submission.problem_id = $2
and public.submission.usr_id = $1
order by public.submission.submission_id desc
limit 1`, uid, pid)

	submission, err := scanSubmissionRow(row)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

// SpecificSubmission return a specific submission by id Ou
func SpecificSubmission(id int) (*Submission, error) {
	row := db.DB.QueryRow(baseSQL + " WHERE submission_id=$1", id)

	submission, err := scanSubmissionRow(row)
	if err != nil {
		return nil, err
	}
	
	return submission, nil
}

// DeleteAllSubmissions cleans all submission
func DeleteAllSubmissions() error {
	statement := `DELETE FROM submission;`
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func collectScore(id int) error {
	statement := 
`update public.grader_user
set score = (
	select sum(max) from (
		select max(public.submission.score * 10.0 / public.submission.max_score) * public.problem.difficulty as max
		from public.submission
		inner join public.problem
		on public.submission.problem_id = public.problem.id
		where public.submission.usr_id = $1
		group by public.submission.problem_id, public.problem.difficulty
	) as score
)
where public.grader_user.id = $1`
	_, err := db.DB.Exec(statement, id)
	if err != nil {
		return err
	}
	return nil
}

func scanSubmissionRow(row *sql.Row) (*Submission, error) {
	submission := new(Submission)

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

func scansSubmissionRow(rows *sql.Rows) ([]*Submission, error) {
	submissions := make([]*Submission, 0)

	defer rows.Close()

	for rows.Next() {
		submission := new(Submission)

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