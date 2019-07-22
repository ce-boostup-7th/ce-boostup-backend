package route

import (
	"../api"
	"net/http"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

//Init init a router for api
func Init() *echo.Echo {
	e := echo.New()

	// config CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1234", "http://problem-injector.surge.sh", "https://ceboostup.netlify.com"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Home
	e.GET("/", api.Home)

	//login&logout route
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
	e.POST("/problems", api.CreateProblem)
	e.GET("/problems", api.GetAllProblems)
	e.GET("/problems/:id", api.GetProblemWithID)
	e.PUT("/problems/:id", api.UpdateProblem)
	e.DELETE("/problems/:id", api.DeleteProblemWithSpecificID)

	e.GET("/problems/:id/testcases", api.GetTestcaseWithID)
	e.POST("/problems/:id/testcases", api.CreateTestcase)
	e.PUT("/problems/:id/testcases/:index", api.UpdateTestcase)
	e.DELETE("/problems/:id/testcases/:index", api.DeleteTestcase)

	//submission routes
	e.POST("/submissions", api.CreateSubmission)
	e.GET("/submissions", api.GetAllSubmissions)
	e.GET("/submissions/:id", api.GetSubmissionWithID)

	e.GET("/users/submissions", api.GetAllSubmissionsOfUser)

	//special
	e.GET("/users/stats", api.GetUserStats)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", api.Restricted)

	return e
}
