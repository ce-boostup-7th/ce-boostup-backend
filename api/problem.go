package api

import (
	"ce-boostup-backend/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

//CreateProblem create a new problem
func CreateProblem(c echo.Context) error {
	values := c.QueryParams()

	//convert string to int
	categoryID, _ := strconv.Atoi(values.Get("categoryID"))

	//convert string to int
	difficulty, _ := strconv.Atoi(values.Get("difficulty"))

	err := model.NewProblem(values.Get("title"), categoryID, difficulty)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusCreated, "a new problem created")
}

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

// GetTestcaseWithID get testcase from judge0
func GetTestcaseWithID(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	testcase, err := model.SpecificTestcaseWithID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "not found any testcases")
	}
	return c.JSON(http.StatusOK, testcase)
}

//CreateTestcase create a new testcase
func CreateTestcase(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, err)
	}

	values := c.QueryParams()

	var testcase model.Testcase
	testcase.Input = values.Get("input")
	testcase.Output = values.Get("output")

	err = model.NewTestcase(id, testcase)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, err)
	}
	return c.JSON(http.StatusOK, "okay")
}

//UpdateProblem update problem data
func UpdateProblem(c echo.Context) error {
	var problem model.Problem

	str := c.Param("id")
	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	problem.ID = id

	values := c.QueryParams()

	if values.Get("title") != "" {
		problem.Title = values.Get("title")
	} else {
		temp, _ := model.SpecificProblemWithID(id)
		problem.Title = temp.Title
	}

	if values.Get("description") != "" {
		problem.Description = values.Get("description")
	} else {
		temp, _ := model.SpecificProblemWithID(id)
		problem.Description = temp.Description
	}

	if values.Get("categoryID") != "" {
		categoryID, _ := strconv.Atoi(values.Get("categoryID"))
		problem.CategoryID = categoryID
	} else {
		temp, _ := model.SpecificProblemWithID(id)
		problem.CategoryID = temp.CategoryID
	}

	if values.Get("difficulty") != "" {
		difficulty, _ := strconv.Atoi(values.Get("difficulty"))
		problem.Difficulty = difficulty
	} else {
		temp, _ := model.SpecificProblemWithID(id)
		problem.Difficulty = temp.Difficulty
	}

	err = model.UpdateProblem(problem)
	if err != nil {
		return c.String(http.StatusNotFound, "update failed")
	}
	return c.String(http.StatusOK, "updated")
}

//DeleteAllProblems delete every problems
func DeleteAllProblems(c echo.Context) error {
	err := model.DeleteAllProblems()
	if err != nil {
		return c.String(http.StatusNotFound, "delete failed")
	}
	return c.String(http.StatusOK, "deleted")
}

//DeleteProblemWithSpecificID delete a problem by id
func DeleteProblemWithSpecificID(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	err1 := model.DeleteProblemWithSpecificID(id)
	if err1 != nil {
		fmt.Println(err1)
		return c.String(http.StatusNotFound, "delete failed")
	}
	return c.String(http.StatusOK, "deleted")
}
