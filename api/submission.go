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

// CreateSubmission create a new submission Ou
func CreateSubmission(c echo.Context) error {
	submission := new(model.Submission)
	if err := c.Bind(submission); err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Request data not in correct format",
			Err: err,
		})
	}

	// read a cookie
	cookie, err := c.Cookie("JWT_Token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Invalid TOKEN",
			Err: err,
		})
	}

	jwtString := cookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Invalid TOKEN",
			Err: err,
		})
	}

	userIDStr := fmt.Sprintf("%v", claims["userID"])
	userID, err := conversion.StringToInt(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	submissionResu, err := model.NewSubmission(userID, submission.ProblemID, submission.LanguageID, submission.Src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Call staff",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated, submissionResu)
}

// GetAllSubmissions get all submissions Ou
func GetAllSubmissions(c echo.Context) error {
	submissions, err := model.AllSubmissions()
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found any problem",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, submissions)
}

// GetAllSubmissionsOfUser get all submissions of specific user Ou
func GetAllSubmissionsOfUser(c echo.Context) error {
	// read a cookie
	cookie, err := c.Cookie("JWT_Token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Invalid TOKEN",
			Err: err,
		})
	}

	jwtString := cookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Invalid TOKEN",
			Err: err,
		})
	}

	userIDStr := fmt.Sprintf("%v", claims["userID"])
	userID, err := conversion.StringToInt(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}


	submissions, err := model.AllSubmissionsFilteredByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found this submissions",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, submissions)
}

// GetSubmissionWithID get a specific submission by id Ou
func GetSubmissionWithID(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	submission, err := model.SpecificSubmission(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found this problem",
			Err: err,
		})
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
