package api

import (
	"ce-boostup-backend/model"
	"log"
	"net/http"
	"strconv"

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

// GetSubmissionWithID get a specific submission by id
func GetSubmissionWithID(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	var submission *model.Submission
	submission, err = model.SpecificSubmission(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, submission)
}
