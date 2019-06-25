package api

import (
	"ce-boostup-backend/model"
	"fmt"
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
	usr, _ := model.AllUsers()
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

//UpdateUser update user data
func UpdateUser(c echo.Context) error {
	var usr model.User

	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	values := c.QueryParams()

	if values.Get("username") != "" {
		usr.Username = values.Get("username")
	} else {
		temp, _ := model.SpecificUserWithID(id)
		usr.Username = temp.Username
	}

	if values.Get("password") != "" {
		usr.Password = values.Get("password")
	} else {
		temp, _ := model.SpecificUserWithID(id)
		usr.Password = temp.Password
	}

	usr.ID = id

	err = model.UpdateUser(usr)
	if err != nil {
		return c.String(http.StatusNotFound, "update failed")
	}
	return c.String(http.StatusOK, "updated")
}

//DeleteAllUsers delete all users
func DeleteAllUsers(c echo.Context) error {
	err := model.DeleteAllUsers()
	if err != nil {
		return c.String(http.StatusNotFound, "delete failed")
	}
	return c.String(http.StatusOK, "deleted")
}

//DeleteUserWithSpecificID delete an user by id
func DeleteUserWithSpecificID(c echo.Context) error {
	str := c.Param("id")

	//convert string to int
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}

	err1 := model.DeleteUserWithSpecificID(id)
	if err1 != nil {
		fmt.Println(err1)
		return c.String(http.StatusNotFound, "delete failed")
	}
	return c.String(http.StatusOK, "deleted")
}
