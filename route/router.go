package route

import (
	"ce-boostup-backend/api"
	"net/http"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

//Init init a router for api
func Init() *echo.Echo {
	e := echo.New()

	// config CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1234"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Home
	e.GET("/", api.Home)

	//login route
	e.POST("/login", api.Login)
	e.POST("/logout", api.Logout)

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

	//submission routes
	e.GET("/submissions", api.GetAllSubmissions)
	e.GET("/submissions/:id", api.GetSubmissionWithID)
	e.POST("/submissions", api.CreateSubmission)
	e.DELETE("/submissions", api.DeleteAllSubmissions)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", api.Restricted)

	return e
}
