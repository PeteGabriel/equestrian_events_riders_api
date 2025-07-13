package application

import (
	"os"

	"github.com/gin-gonic/gin"
)

func (app *Application) Serve() error {

	r := app.routes()

	// grab port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return r.Run(":" + port)
}

func (app *Application) routes() *gin.Engine {
	router := gin.Default()

	router.
		GET("/", app.HandleCompetitions(app.ListCompetitions)).
		GET("/competitions", app.HandleCompetitions(app.ListCompetitions))

	router.GET("/competitions/:id", app.HandleCompetitionByID(app.GetCompetitionByID))

	return router
}
