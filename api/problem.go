package api

import (
	"../conversion"
	"../model"
	"net/http"

	"github.com/labstack/echo"
)

//RespSuccess struct for json return
type RespSuccess struct {
  	ID int `json:"id"`
	Msg string `json:"msg"`
}

//RespError struct for json return
type RespError struct {
  Msg string `json:"msg"`
  Err error `json:"err"`
}

// CreateProblem create a new problem Ou
func CreateProblem(c echo.Context) error {
	problem := new(model.Problem)
	if err := c.Bind(problem); err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Request data not in correct format",
			Err: err,
		})
	}

	id, err := model.NewProblem(problem.Title, problem.CategoryID, problem.Difficulty, problem.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Can not create problem",
			Err: err,
		})
	}

	return c.JSON(http.StatusCreated, &RespSuccess{ID: *id, Msg: "Created"})
}

// GetAllProblems get all problems Ou
func GetAllProblems(c echo.Context) error {
	problems, err := model.AllProblems()
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found any problem",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, problems)
}

// GetProblemWithID get specific problem by id Ou
func GetProblemWithID(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	problem, err := model.SpecificProblemWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found that problem ID",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, problem)
}

// UpdateProblem update problem data Ou
func UpdateProblem(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	problem, err := model.SpecificProblemWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Can not found that problem ID",
			Err: err,
		})
	}
	if err = c.Bind(problem); err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "Request data not in correct format",
			Err: err,
		})
	}

	err = model.UpdateProblem(*problem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Can not update that problem ID",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, &RespSuccess{ID: id, Msg: "Updated"})
}

//DeleteAllProblems delete every problems
func DeleteAllProblems(c echo.Context) error {
	err := model.DeleteAllProblems()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}

// DeleteProblemWithSpecificID delete a problem by id Ou
func DeleteProblemWithSpecificID(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	err = model.DeleteProblemWithSpecificID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Can not delete that problem ID",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK,  &RespSuccess{ID: id, Msg: "Delete"})
}

// CreateTestcase create a new testcase Ou
func CreateTestcase(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	testcase := new(model.Testcase)

	if err := c.Bind(testcase); err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg:  "Request data not in correct format",
			Err: err,
		})
	}
	err = model.NewTestcase(id, *testcase)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg:  "Can not create new testcase",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated,  &RespSuccess{Msg: "Created"})
}

// GetTestcaseWithID get testcase from judge0 Ou
func GetTestcaseWithID(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	testcase, err := model.SpecificTestcaseWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found any testcase",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, testcase)
}


// UpdateTestcase create a new testcase Ou
func UpdateTestcase(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	index, err := conversion.StringToInt(c.Param("index"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	index++;

	testcase := new(model.Testcase)

	if err := c.Bind(testcase); err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg:  "Request data not in correct format",
			Err: err,
		})
	}
	err = model.UpdateTestcase(id, index, *testcase)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg:  "Can not create new testcase",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated,  &RespSuccess{Msg: "Update"})
}

// DeleteTestcase create a new testcase Ou
func DeleteTestcase(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	index, err := conversion.StringToInt(c.Param("index"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	index++;

	err = model.DeleteTestcase(id, index)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Can not delete that problem ID",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated,  &RespSuccess{Msg: "Deleted"})
}