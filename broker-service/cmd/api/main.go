package main

import (
	"log"
	"net/http"
)

const port = ":8000"

type application struct{}

func main() {
	app := application{}
	log.Printf("Starting broker on port %s\n", port)
	srv := &http.Server{
		Addr:    port,
		Handler: app.routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
