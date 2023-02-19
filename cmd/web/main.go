package main

import (
	"flag"
	"kvstore/kv"
	"kvstore/kv/memory"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct {
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	Store           kv.Store
	MaxPayloadBytes int64
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	store := memory.NewStore()

	app := &Application{
		InfoLog:         infoLog,
		ErrorLog:        errorLog,
		Store:           store,
		MaxPayloadBytes: 1024 * 5, // 5 KiB payload
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
