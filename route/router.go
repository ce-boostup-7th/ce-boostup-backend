package route

import (
	"ce-boostup-backend/api"

	"github.com/labstack/echo"
)

//Init init a router for api
func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", api.Home)
	e.GET("/user", api.GetUser)

	return e
}
