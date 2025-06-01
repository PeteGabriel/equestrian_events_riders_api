package application

import (
	"github.com/gin-gonic/gin"
	"os"
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

	router.Use(app.CheckCacheForEntryLists).
		GET("/competitions", app.HandleCompetitions(app.ListCompetitions))

	return router
}
