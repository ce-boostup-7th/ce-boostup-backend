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

	e.GET("/users/submissions", api.GetAllSubmissionsOfUser)

	//special
	e.GET("/users/stats", api.GetUserStats)

	// ---------- OU version ----------

	// login&logout route
	e.POST("/ou/login", api.OuLogin)
	e.POST("/ou/logout", api.OuLogout)

	//problem routes
	e.POST("/ou/problems", api.OuCreateProblem)
	e.GET("/ou/problems", api.OuGetAllProblems)
	e.GET("/ou/problems/:id", api.OuGetProblemWithID)
	e.PUT("/ou/problems/:id", api.OuUpdateProblem)
	e.DELETE("/ou/problems/:id", api.OuDeleteProblemWithSpecificID)

	e.GET("/ou/problems/:id/testcases", api.OuGetTestcaseWithID)
	e.POST("/ou/problems/:id/testcases", api.OuCreateTestcase)
	e.PUT("/ou/problems/:id/testcases/:index", api.OuUpdateTestcase)
	e.DELETE("/ou/problems/:id/testcases/:index", api.OuDeleteTestcase)

	//submission routes
	e.POST("/ou/submissions", api.OuCreateSubmission)
	e.GET("/ou/submissions", api.OuGetAllSubmissions)
	e.GET("/ou/submissions/:id", api.OuGetSubmissionWithID)

	e.GET("/ou/users/submissions", api.OuGetAllSubmissionsOfUser)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", api.Restricted)

	return e
}
