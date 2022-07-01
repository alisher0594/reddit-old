package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	router.HandlerFunc(http.MethodPost, "/v1/posts", app.createPostHandler)
	router.HandlerFunc(http.MethodGet, "/v1/posts/:id", app.showPostHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/posts/:id", app.updatePostHandler)
	router.HandlerFunc(http.MethodGet, "/v1/posts", app.listPostsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/posts/:id", app.deletePostHandler)

	router.HandlerFunc(http.MethodPost, "/v1/posts/:id/up", app.upPostHandler)
	router.HandlerFunc(http.MethodPost, "/v1/posts/:id/down", app.downPostHandler)

	return app.metrics(app.recoverPanic(app.enableCORS(router)))
}
