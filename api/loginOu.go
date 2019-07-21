package api

import (
	"../model"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// ---------- OU version ----------

// OuLogin authorize and return a cookie
func OuLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	hashedPassword, err := model.PasswordByUsername(username)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Incorrect Username or Password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(password))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Incorrect Username or Password")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	userID, err := model.IDByUsername(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "contact admin")
	}

	//Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	//Generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.String(http.StatusInternalServerError, "contact admin")
	}

	cookie := new(http.Cookie)
	cookie.HttpOnly = false
	cookie.Name = "JWT_Token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "logged in")
}