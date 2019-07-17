package api

import (
	"../conversion"
	"../model"
	"net/http"

	"github.com/labstack/echo"
)

//CreateProblem create a new problem
func CreateProblem(c echo.Context) error {
	var problem model.Problem
	if err := c.Bind(&problem); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	model.NewProblem(problem.Title, problem.CategoryID, problem.Difficulty, problem.Description)
	return c.JSON(http.StatusCreated, "created")
}

//GetAllProblems get all problems
func GetAllProblems(c echo.Context) error {
	problems, err := model.AllProblems()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, problems)
}

//GetProblemWithID get specific problem by id
func GetProblemWithID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	problem, err := model.SpecificProblemWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, problem)
}

// GetTestcaseWithID get testcase from judge0
func GetTestcaseWithID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	testcase, err := model.SpecificTestcaseWithID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "not found any testcases")
	}
	return c.JSON(http.StatusOK, testcase)
}

//CreateTestcase create a new testcase
func CreateTestcase(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	var testcase model.Testcase
	if err := c.Bind(&testcase); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	model.NewTestcase(id, testcase)
	return c.JSON(http.StatusCreated, "created")
}

//UpdateProblem update problem data
func UpdateProblem(c echo.Context) error {

	str := c.Param("id")
	id := conversion.StringToInt(str)

	problemPtr, _ := model.SpecificProblemWithID(id)
	problem := *problemPtr
	if err := c.Bind(&problem); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := model.UpdateProblem(problem)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "updated")
}

//DeleteAllProblems delete every problems
func DeleteAllProblems(c echo.Context) error {
	err := model.DeleteAllProblems()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}

//DeleteProblemWithSpecificID delete a problem by id
func DeleteProblemWithSpecificID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	err := model.DeleteProblemWithSpecificID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}
