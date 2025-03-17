package main

import "os"

func (a *Application) serve() error {

	r := a.routes()

	// grab port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return r.Run(":" + port)
}
