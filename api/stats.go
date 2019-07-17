package api

import (
	"ce-boostup-backend/conversion"
	"ce-boostup-backend/model"
	"net/http"

	"github.com/labstack/echo"
)

//GetUserStats get user statistics
func GetUserStats(c echo.Context) error {
	str := c.Param("id")

	id := conversion.StringToInt(str)

	stat, err := model.SpecificUserStatWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, stat)
}
