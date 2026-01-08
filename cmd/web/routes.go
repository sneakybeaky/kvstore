package main

import (
	"net/http"
)

func (app *Application) Routes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("GET /ping", app.logRequest(timeRequest(http.HandlerFunc(ping))))
	mux.Handle("PUT /store/{key}", app.logRequest(timeRequest(http.HandlerFunc(app.set))))
	mux.Handle("GET /store/{key}", app.logRequest(timeRequest(http.HandlerFunc(app.get))))

	return mux
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
