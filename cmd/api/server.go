package main

func (a *Application) serve() error {

	r := a.routes()

	return r.Run("localhost:8084")
}
