package api

import (
	"ce-boostup-backend/model"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//Login authorize and return a cookie
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	//Throws unauthorized error
	if !isPasswordCorrect(username, password) {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	//Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "CE 1D"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	//Generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "JWT"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "logged in")
}

//Accessible accessible without authentication
func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

//Restricted cannot access without authentication
func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func isPasswordCorrect(username string, password string) bool {
	hashedPassword, _ := model.PasswordByUsername(username)
	err := bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(password))
	return err == nil
}

func writeCookie(c echo.Context) {
}

func readCookie(c echo.Context) {
	cookie, _ := c.Cookie("username")
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
}
