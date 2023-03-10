package main

import (
	"fmt"
	"log"
	"net/http"
)

var webPort = "80"

type Config struct {
}

func main() {
	log.Printf("broker_service is running on port %s", webPort)

	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
