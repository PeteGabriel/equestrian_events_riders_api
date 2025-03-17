package main

import "github.com/gin-gonic/gin"

func (a *Application) routes() *gin.Engine {
	router := gin.Default()
	router.Use(a.CheckCacheForEntryLists).Handle("GET", "/competitions", a.ListCompetitions)
	return router
}
