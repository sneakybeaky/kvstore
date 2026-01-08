package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

type value struct {
	Value string
}

func (app *Application) set(w http.ResponseWriter, r *http.Request) {

	key := r.PathValue("key")

	// Use http.MaxBytesReader to enforce a maximum read from the
	// response body.
	r.Body = http.MaxBytesReader(w, r.Body, app.MaxPayloadBytes)

	dec := json.NewDecoder(r.Body)

	var v value
	err := dec.Decode(&v)

	if err != nil {
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &maxBytesError):
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", app.MaxPayloadBytes)
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		default:
			app.ErrorLog.Printf("Unable to decode value : %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}

		return
	}

	app.Store.Set(key, v.Value)
}

func (app *Application) get(w http.ResponseWriter, r *http.Request) {

	key := r.PathValue("key")

	found, ok := app.Store.Get(key)

	if !ok {
		http.Error(w, "no value set for key "+key, http.StatusNotFound)
		return
	}

	body, err := json.Marshal(value{Value: found})

	if err != nil {
		app.ErrorLog.Printf("Unable to encode value : %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}
