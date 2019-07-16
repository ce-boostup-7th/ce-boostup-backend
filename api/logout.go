package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Logout log out
func Logout(c echo.Context) error {

	cookie := new(http.Cookie)
	cookie.Name = "JWT_Token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(1 * time.Microsecond)
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "logged out")
}
