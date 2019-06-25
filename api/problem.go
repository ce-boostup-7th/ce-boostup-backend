package api

import (
	"ce-boostup-backend/model"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

//GetAllProblems get all problems
func GetAllProblems(c echo.Context) error {
	problems, _ := model.AllProblems()
	return c.JSON(http.StatusOK, problems)
}

//GetProblemWithID get specific problem by id
func GetProblemWithID(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	var problem *model.Problem
	problem, err = model.SpecificProblemWithID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "not found that problem")
	}
	return c.JSON(http.StatusOK, problem)
}
