package model

import (
	"ce-boostup-backend/conversion"
	"ce-boostup-backend/db"
	"ce-boostup-backend/judge0"
	"fmt"
)

//Submission a model for submission
type Submission struct {
	SubmissionID int     `json:"submission_id" form:"submission_id"`
	UserID       int     `json:"user_id" form:"user_id"`
	ProblemID    int     `json:"problem_id" form:"problem_id"`
	LanguageID   int     `json:"language_id" form:"language_id"`
	Src          string  `json:"src" form:"src"`
	SubmittedAt  string  `json:"submitted_at" form:"submitted_at"`
	Score        int     `json:"score" form:"score"`
	MaxScore     int     `json:"max_score" form:"max_score"`
	Runtime      float64 `json:"runtime" form:"runtime"`
	MemoryUsage  int     `json:"memory_usage" form:"memory_usage"`
}

//NewSubmission create a new submission
func NewSubmission(userID int, problemID int, languageID int, src string) error {

	score := 0
	runtime := 0.0
	memory := 0

	testcase, err := SpecificTestcaseWithID(problemID)
	if err != nil {
		return err
	}

	for i := range testcase {
		result := judge0.Submit(src, testcase[i].Input, testcase[i].Output) //empty string is for testcase in the future
		memory += result.Memory
		runtime += conversion.StringToFloat(result.Time)
		if result.Status.ID == 3 {
			score++
		}
	}

	length := len(testcase)
	runtime = runtime / float64(length)
	memory = memory / length

	statement := `INSERT INTO submission (usr_id,problem_id,lang_id,src,score,runtime,memory_usage,max_score) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err = db.DB.Exec(statement, userID, problemID, languageID, src, score, runtime, memory, length)
	if err != nil {
		return err
	}
	collectScore(userID)

	return nil
}

// AllSubmissions get all submissions
func AllSubmissions() ([]*Submission, error) {
	rows, err := db.DB.Query("SELECT submission_id,src,usr_id,problem_id,lang_id,submittedat,score,max_score,runtime,memory_usage FROM submission ORDER BY submission_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*Submission, 0)

	for rows.Next() {
		submission := new(Submission)

		err := rows.Scan(&submission.SubmissionID, &submission.Src, &submission.UserID, &submission.ProblemID, &submission.LanguageID, &submission.SubmittedAt, &submission.Score, &submission.MaxScore, &submission.Runtime, &submission.MemoryUsage)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

//AllSubmissionsFilteredByUserID filter all user that filtered by userID
func AllSubmissionsFilteredByUserID(uid int) ([]*Submission, error) {
	rows, err := db.DB.Query("SELECT submission_id,src,usr_id,problem_id,lang_id,submittedat,score,max_score,runtime,memory_usage FROM submission WHERE usr_id=$1 ORDER BY submission_id", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*Submission, 0)

	for rows.Next() {
		submission := new(Submission)

		err := rows.Scan(&submission.SubmissionID, &submission.Src, &submission.UserID, &submission.ProblemID, &submission.LanguageID, &submission.SubmittedAt, &submission.Score, &submission.MaxScore, &submission.Runtime, &submission.MemoryUsage)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

// SpecificSubmission return a specific submission by id
func SpecificSubmission(id int) (*Submission, error) {
	statement := `SELECT submission_id,src,usr_id,problem_id,lang_id,submittedat,score,runtime,memory_usage  FROM submission WHERE submission_id=$1`
	row := db.DB.QueryRow(statement, id)

	submission := new(Submission)

	err := row.Scan(&submission.SubmissionID, &submission.Src, &submission.UserID, &submission.LanguageID, &submission.LanguageID, &submission.SubmittedAt, &submission.Score, &submission.Runtime, &submission.MemoryUsage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return submission, nil
}

// DeleteAllSubmissions cleans all submission
func DeleteAllSubmissions() error {
	statement := `DELETE FROM submission; ALTER SEQUENCE submission_submission_id_seq RESTART WITH 1;`
	_, err := db.DB.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func collectScore(id int) error {
	statement := `UPDATE grader_user SET score=(SELECT SUM(max) FROM (SELECT problem_id,MAX(score) FROM submission WHERE usr_id=$1 GROUP BY submission.problem_id) AS PREP) WHERE id=$1;`
	_, err := db.DB.Exec(statement, id)
	if err != nil {
		return err
	}
	return nil
}
