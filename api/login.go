package api

import (
	"ce-boostup-backend/model"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//RespSuccess struct for json return
type RespSuccess struct {
	ID  int    `json:"id"`
	Msg string `json:"msg"`
}

//RespError struct for json return
type RespError struct {
	Msg string `json:"msg"`
	Err error  `json:"err"`
}

// Login authorize and return a cookie Ou
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	isExist, _ := model.IsUserExist(username)
	//Throws unanthorized error
	if !(*isExist) {
		return c.String(http.StatusNotFound, "cannot found that user")
	}

	userID, hashedPassword, err := model.IDPasswordByUsername(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "CallStaff")
	}
	if *userID == -1 {
		return c.String(http.StatusNotFound, "NoUser")
	}
	err = bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(password))
	if err != nil {
		return c.String(http.StatusUnauthorized, "IncorrectPassword")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	endTime := time.Now().Add(time.Hour * 24)

	//Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = endTime.Unix()

	//Generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.String(http.StatusInternalServerError, "contact admin")
	}

	cookie := new(http.Cookie)
	cookie.Name = "JWT_Token"
	cookie.Value = t
	cookie.Expires = endTime
	cookie.Path = "/"
	cookie.HttpOnly = false
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "logged in")
}
