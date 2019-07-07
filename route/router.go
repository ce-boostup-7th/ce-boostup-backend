package route

import (
	"ce-boostup-backend/api"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

//Init init a router for api
func Init() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Home
	e.GET("/", api.Home)

	//login route
	e.POST("/login", api.Login)

	//user routes
	e.GET("/users", api.GetAllUsers)
	e.GET("/users/:id", api.GetUserWithID)
	e.POST("/users", api.CreateUser)
	e.PUT("/users/:id", api.UpdateUser)
	e.DELETE("/users", api.DeleteAllUsers)
	e.DELETE("/users/:id", api.DeleteUserWithSpecificID)

	//problem routes
	e.GET("/problems", api.GetAllProblems)
	e.GET("/problems/:id", api.GetProblemWithID)
	e.POST("/problems", api.CreateProblem)
	e.PUT("/problems/:id", api.UpdateProblem)
	e.DELETE("/problems", api.DeleteAllProblems)
	e.DELETE("problems/:id", api.DeleteProblemWithSpecificID)
	e.GET("/problems/testcases/:id", api.GetTestcaseWithID)
	e.POST("/problems/testcases/:id", api.CreateTestcase)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", api.Restricted)

	return e
}
