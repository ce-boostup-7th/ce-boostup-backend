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
		AllowOrigins:     []string{"http://localhost:1234", "http://problem-injector.surge.sh", "https://ceboostup.netlify.com", "http://boostup-demo.surge.sh"},
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

	//problem routes
	e.GET("/problems", api.GetAllProblems)
	e.GET("/problems/:id", api.GetProblemWithID)
	e.GET("/problems/:id/testcases", api.GetTestcaseWithID)
	
	//submission routes
	e.POST("/submissions", api.CreateSubmission)
	e.GET("/submissions", api.GetAllSubmissions)
	e.GET("/submissions/:id", api.GetSubmissionWithID)

	e.GET("/users/submissions", api.GetAllSubmissionsOfUser)
	e.GET("/users/problems", api.GetAllProblemsWithUserProgres)
	e.GET("/problems/:pid/submissions/last", api.GetLastUserSubmissionsFilteredByProblemID)

	// ---------- only admin ----------
	r := e.Group("/Ad-0Fj_kL8me")

	// User routes
	r.PUT("/users/:id", api.UpdateUser)
	r.DELETE("/users/:id", api.DeleteUserWithSpecificID)

	// problem routes
	r.POST("/problems", api.CreateProblem)
	r.PUT("/problems/:id", api.UpdateProblem)
	r.DELETE("/problems/:id", api.DeleteProblemWithSpecificID)

	r.POST("/problems/:id/testcases", api.CreateTestcase)
	r.GET("/problems/:id/testcases", api.GetTestcaseWithIDAll)
	r.PUT("/problems/:id/testcases/:index", api.UpdateTestcase)
	r.DELETE("/problems/:id/testcases/:index", api.DeleteTestcase)

	//special
	e.GET("/users/stats", api.GetUserStats)

	return e
}
