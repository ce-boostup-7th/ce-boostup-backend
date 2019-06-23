package api

import (
	"ce-boostup-backend/model"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

//CreateUser create a new user
func CreateUser(c echo.Context) error {
	values := c.QueryParams()

	err := model.NewUser(values.Get("username"), values.Get("password"))
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusCreated, "a new user created")
}

//GetAllUsers get all users info
func GetAllUsers(c echo.Context) error {
	var usr []*model.User

	usr, _ = model.AllUsers()
	return c.JSON(http.StatusOK, usr)
}

//GetUserWithID Get specific user with id
func GetUserWithID(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	var usr *model.User
	usr, err = model.SpecificUserWithID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "not found that user")
	}
	return c.JSON(http.StatusOK, usr)
}
