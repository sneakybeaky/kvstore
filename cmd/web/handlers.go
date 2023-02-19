package main

import (
	"encoding/json"
	"net/http"
)

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

type keyValue struct {
	Key   string
	Value string
}

func (app *Application) store(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)

	var kv keyValue
	err := dec.Decode(&kv)

	if err != nil {

		app.ErrorLog.Printf("Unable to decode store request body : %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadGateway)
		return
	}

	app.Store.Set(kv.Key, kv.Value)
}
