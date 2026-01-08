package main

import (
	"log/slog"
	"net/http"
	"time"
)

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)

	})
}

func timeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			start = time.Now()
		)

		next.ServeHTTP(w, r)
		slog.Info("request time", slog.Duration("duration", time.Since(start)))
	})
}
