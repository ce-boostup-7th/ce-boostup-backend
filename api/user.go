package api

import (
	"ce-boostup-backend/model"
	"net/http"

	"github.com/labstack/echo"
)

//GetUser get users info
func GetUser(c echo.Context) error {
	var usr []*model.User

	usr, _ = model.AllUsers()
	return c.JSON(http.StatusOK, usr)
}
