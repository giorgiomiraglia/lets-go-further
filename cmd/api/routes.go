package main

import "github.com/julienschmidt/httprouter"

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc("GET", "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc("GET", "/v1/movies/:id", app.showMoviesHandler)

	router.HandlerFunc("POST", "/v1/movies", app.createMoviesHandler)

	return router
}
