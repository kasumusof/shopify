package main

import (
	"log"
	"net/http"

	"github.com/kasumusof/shopify/pkg"
)

func init() {
	log.Println("Starting server...")
	log.Println("Listen on port 8080")
	log.Println("Press CTRL+C to stop")
}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: pkg.Router(),
	}

	server.ListenAndServe()
}
