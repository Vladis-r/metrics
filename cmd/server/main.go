package main

import (
	"net/http"

	"github.com/Vladis-r/metrics.git/internal/handler"
)

func main() {
	// Create custom mux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/", handler.MainHandler)
	mux.HandleFunc("/update/", handler.UpdateHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
