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

	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//hash a password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	err := model.NewUser(user.Username, string(bytes))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.String(http.StatusCreated, "a new user created")
}

//GetAllUsers get all users info
func GetAllUsers(c echo.Context) error {
	user, err := model.AllUsers()
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, user)
}

//GetUserWithID Get specific user with id
func GetUserWithID(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	user, err := model.SpecificUserWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, user)
}

//UpdateUser update user data
func UpdateUser(c echo.Context) error {

	str := c.Param("id")
	id := conversion.StringToInt(str)

	userPtr, _ := model.SpecificUserWithID(id)
	user := *userPtr

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := model.UpdateUser(user)
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
