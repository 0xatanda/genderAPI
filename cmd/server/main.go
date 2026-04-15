package main

import (
	"net/http"

	"github.com/0xatanda/genderAPI/internal/handler"
)

func main() {
	http.HandleFunc("/api/classify", handler.ClassifyHandler)
	http.ListenAndServe(":80", nil)
}
