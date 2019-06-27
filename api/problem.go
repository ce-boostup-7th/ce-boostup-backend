package api

import (
	"ce-boostup-backend/model"
	"fmt"
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
