package api

import (
	"ce-boostup-backend/model"
	"net/http"

	"github.com/labstack/echo"
)

// CreateSubmission create a new submission
func CreateSubmission(c echo.Context) error {
	var submission model.Submission
	if err := c.Bind(&submission); err != nil {
		return err
	}
	model.NewSubmission(submission.UserID, submission.ProblemID, submission.LanguageID, submission.Src)
	return c.JSON(http.StatusCreated, "created")
}

// GetAllSubmissions get all submissions
func GetAllSubmissions(c echo.Context) error {
	submissions, err := model.AllSubmissions()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, submissions)
}
