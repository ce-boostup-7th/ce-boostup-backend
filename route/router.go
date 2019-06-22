package route

import (
	"ce-boostup-backend/api"

	"github.com/labstack/echo"
)

//Init init a router for api
func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", api.Home)

	//user handlers
	e.GET("/user", api.GetAllUsers)
	e.GET("/user/:id", api.GetUserWithID)

	return e
}
