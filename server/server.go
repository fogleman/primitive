// main.go
package main

import (
	"log"
	"net/http"
	"os"

	pr "github.com/fogleman/primitive/server/route"
)

func setupRoutes() {
	handler := pr.PrimitiveRoute(pr.Config{
		MaxUploadMb: 10,
	})
	http.HandleFunc("/primitive", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("We're running!"))
	})
}

func main() {
	setupRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server - ", err)
	}
}
