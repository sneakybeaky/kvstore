package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

type value struct {
	Value string
}

func (app *Application) store(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	key := params.ByName("key")

	dec := json.NewDecoder(r.Body)

	var v value
	err := dec.Decode(&v)

	if err != nil {

		app.ErrorLog.Printf("Unable to decode value : %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadGateway)
		return
	}

	app.Store.Set(key, v.Value)
}
