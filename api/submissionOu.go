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

// ---------- OU version ----------

// OuCreateSubmission create a new submission
func OuCreateSubmission(c echo.Context) error {
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

	id, err := model.OuNewSubmission(userID, submission.ProblemID, submission.LanguageID, submission.Src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Call staff",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated, &RespSuccess{ID: id, Msg: "Judge success"})
}

// OuGetAllSubmissions get all submissions
func OuGetAllSubmissions(c echo.Context) error {
	submissions, err := model.OuAllSubmissions()
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found any problem",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, submissions)
}

// OuGetAllSubmissionsOfUser get all submissions of specific user
func OuGetAllSubmissionsOfUser(c echo.Context) error {
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


	submissions, err := model.OuAllSubmissionsFilteredByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found this submissions",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, submissions)
}

// OuGetSubmissionWithID get a specific submission by id
func OuGetSubmissionWithID(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	submission, err := model.OuSpecificSubmission(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found this problem",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, submission)
}