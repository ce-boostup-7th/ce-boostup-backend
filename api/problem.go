package api

import (
	"ce-boostup-backend/model"
	"net/http"

	"github.com/labstack/echo"
)

//GetAllProblems get all problems
func GetAllProblems(c echo.Context) error {
	problems, _ := model.AllProblems()
	return c.JSON(http.StatusOK, problems)
}
