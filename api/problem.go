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
	id, err := model.NewProblem(problem.Title, problem.CategoryID, problem.Difficulty, problem.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, id)
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

	id, _ := conversion.StringToInt(str)

	problem, err := model.SpecificProblemWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, problem)
}

// GetTestcaseWithID get testcase from judge0
func GetTestcaseWithID(c echo.Context) error {
	str := c.Param("id")

	id, _ := conversion.StringToInt(str)

	testcase, err := model.SpecificTestcaseWithID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "not found any testcases")
	}
	return c.JSON(http.StatusOK, testcase)
}

//CreateTestcase create a new testcase
func CreateTestcase(c echo.Context) error {
	str := c.Param("id")

	id, _ := conversion.StringToInt(str)

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
	id, _ := conversion.StringToInt(str)

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

	id, _ := conversion.StringToInt(str)

	err := model.DeleteProblemWithSpecificID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}

// ---------- OU version ----------

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

// OuCreateProblem create a new problem
func OuCreateProblem(c echo.Context) error {
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

// OuGetAllProblems get all problems
func OuGetAllProblems(c echo.Context) error {
	problems, err := model.AllProblems()
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found any problem",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, problems)
}

// OuGetProblemWithID get specific problem by id
func OuGetProblemWithID(c echo.Context) error {
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


// OuUpdateProblem update problem data
func OuUpdateProblem(c echo.Context) error {
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

// OuDeleteProblemWithSpecificID delete a problem by id
func OuDeleteProblemWithSpecificID(c echo.Context) error {
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

// OuCreateTestcase create a new testcase
func OuCreateTestcase(c echo.Context) error {
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
	err = model.OuNewTestcase(id, *testcase)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg:  "Can not create new testcase",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated,  &RespSuccess{Msg: "Created"})
}

// OuGetTestcaseWithID get testcase from judge0
func OuGetTestcaseWithID(c echo.Context) error {
	id, err := conversion.StringToInt(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &RespError{
			Msg: "ID can only be integer",
			Err: err,
		})
	}

	testcase, err := model.OuSpecificTestcaseWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &RespError{
			Msg: "Not found any testcase",
			Err: err,
		})
	}
	return c.JSON(http.StatusOK, testcase)
}

// OuUpdateTestcase create a new testcase
func OuUpdateTestcase(c echo.Context) error {
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
	err = model.OuUpdateTestcase(id, index, *testcase)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg:  "Can not create new testcase",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated,  &RespSuccess{Msg: "Update"})
}

// OuDeleteTestcase create a new testcase
func OuDeleteTestcase(c echo.Context) error {
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

	err = model.OuDeleteTestcase(id, index)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &RespError{
			Msg: "Can not delete that problem ID",
			Err: err,
		})
	}
	return c.JSON(http.StatusCreated,  &RespSuccess{Msg: "Deleted"})
}