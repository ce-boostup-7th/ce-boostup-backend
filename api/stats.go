package api

import (
	"../conversion"
	"../model"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//GetUserStats get user statistics
func GetUserStats(c echo.Context) error {
	// read a cookie
	cookie, err := c.Cookie("JWT_Token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	jwtString := cookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userIDStr := fmt.Sprintf("%v", claims["userID"])
	id, _ := conversion.StringToInt(userIDStr)
	stat, err := model.SpecificUserStatWithID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, stat)
}
