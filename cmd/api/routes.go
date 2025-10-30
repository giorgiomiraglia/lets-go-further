package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc("GET", "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc("GET", "/v1/movies/:id", app.showMoviesHandler)
	router.HandlerFunc("POST", "/v1/movies", app.createMoviesHandler)

	return router
}
