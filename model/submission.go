package model

import (
	"ce-boostup-backend/db"
)

//Submission a model for submission
type Submission struct {
	SubmissionID int
	UserID       int
	ProblemID    int
	LanguageID   int
	Src          string
	SubmittedAt  string
	Score        int
	Runtime      int
	MemoryUsage  float32
}

//NewSubmission create a new submission
func NewSubmission(userID int, problemID int, languageID int, src string) error {

	score := 0
	runtime := 0.0
	memory := 0

	//result := judge0.Submit(src, "", "") //empty string is for testcase in the future

	statement := `INSERT INTO submission (usr_id,problem_id,lang_id,src,score,runtime,memory_usage) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := db.DB.Exec(statement, userID, problemID, languageID, src, score, runtime, memory)
	if err != nil {
		return err
	}
	return nil
	//wait for connecting to Judge0 api
}
