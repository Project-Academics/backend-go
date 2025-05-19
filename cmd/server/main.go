package main

import (
	"log"
	"net/http"

	"github.com/Project-Academics/backend-go/internal/upload"
)

func main() {
	http.HandleFunc("/api/upload", upload.Handler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
