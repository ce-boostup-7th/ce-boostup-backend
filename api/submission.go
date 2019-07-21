package api

import (
	"../conversion"
	"../model"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

// CreateSubmission create a new submission
func CreateSubmission(c echo.Context) error {
	var submission model.Submission
	if err := c.Bind(&submission); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// read a cookie
	cookie, err := c.Cookie("JWT_Token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	jwtString := cookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userIDStr := fmt.Sprintf("%v", claims["userID"])
	userID := conversion.StringToInt(userIDStr)

	model.NewSubmission(userID, submission.ProblemID, submission.LanguageID, submission.Src)
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

//GetAllSubmissionsOfUser get all submissions of specific user
func GetAllSubmissionsOfUser(c echo.Context) error {
	// read a cookie
	cookie, err := c.Cookie("JWT_Token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	jwtString := cookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userIDStr := fmt.Sprintf("%v", claims["userID"])
	userID := conversion.StringToInt(userIDStr)

	problems, err := model.AllSubmissionsFilteredByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, problems)
}

// GetSubmissionWithID get a specific submission by id
func GetSubmissionWithID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	var submission *model.Submission
	submission, err := model.SpecificSubmission(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, submission)
}

// DeleteAllSubmissions delete all submissions
func DeleteAllSubmissions(c echo.Context) error {
	err := model.DeleteAllSubmissions()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}
