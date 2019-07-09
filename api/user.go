package api

import (
	"ce-boostup-backend/conversion"
	"ce-boostup-backend/model"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser create a new user
func CreateUser(c echo.Context) error {
	values := c.QueryParams()

	//hash a password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(values.Get("password")), 14)

	err := model.NewUser(values.Get("username"), string(bytes))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.String(http.StatusCreated, "a new user created")
}

//GetAllUsers get all users info
func GetAllUsers(c echo.Context) error {
	usr, err := model.AllUsers()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, usr)
}

//GetUserWithID Get specific user with id
func GetUserWithID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	usr, err := model.SpecificUserWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, usr)
}

//UpdateUser update user data
func UpdateUser(c echo.Context) error {
	var usr model.User

	str := c.Param("id")
	id := conversion.StringToInt(str)
	usr.ID = id

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

	err := model.UpdateUser(usr)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "updated")
}

//DeleteAllUsers delete all users
func DeleteAllUsers(c echo.Context) error {
	err := model.DeleteAllUsers()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}

//DeleteUserWithSpecificID delete an user by id
func DeleteUserWithSpecificID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	err := model.DeleteUserWithSpecificID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}
