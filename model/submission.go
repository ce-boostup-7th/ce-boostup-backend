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
	statement := `INSERT INTO submission (usr_id,problem_id,lang_id,src) VALUES ($1,$2,$3,$4)`
	_, err := db.DB.Exec(statement, userID, problemID, languageID, src)
	if err != nil {
		return err
	}
	return nil
	//wait for connecting to Judge0 api
}
