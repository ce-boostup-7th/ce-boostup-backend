package api

import (
	"ce-boostup-backend/conversion"
	"ce-boostup-backend/model"
	"net/http"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
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
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

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

	id, _ := conversion.StringToInt(str)

	user, err := model.SpecificUserWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, user)
}

//UpdateUser update user data
func UpdateUser(c echo.Context) error {

	str := c.Param("id")
	id, _ := conversion.StringToInt(str)

	userPtr, _ := model.SpecificUserWithID(id)
	user := new(model.User)
	user.ID = userPtr.ID
	user.Username = userPtr.Username
	user.Score = userPtr.Score

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//hash a password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)

	err := model.UpdateUser(*user)
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

	id, _ := conversion.StringToInt(str)

	err := model.DeleteUserWithSpecificID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.String(http.StatusOK, "deleted")
}

func getUserID(c echo.Context) (int, error) {
	// read a cookie
	cookie, err := c.Cookie("JWT_Token")
	if err != nil {
		return -1, err
	}

	jwtString := cookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return -1, err
	}

	userIDStr := fmt.Sprintf("%v", claims["userID"])
	userID, err := conversion.StringToInt(userIDStr)
	if err != nil {
		return -1, err
	}

	return userID, nil
}