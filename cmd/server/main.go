package main

import (
	"log"
	"net/http"
	"os"

	"github.com/0xatanda/genderAPI/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/classify", handler.ClassifyHandler)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
