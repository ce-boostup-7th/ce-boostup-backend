package api

import (
	"net/http"

	"github.com/labstack/echo"
)

//Home home for api
func Home(context echo.Context) error {
	return context.String(http.StatusOK, "Hello, CE Boostup")
}
