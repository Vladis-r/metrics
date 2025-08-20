package main

import (
	"encoding/json"
	"net/http"

	"github.com/Vladis-r/metrics.git/internal/handler"
	models "github.com/Vladis-r/metrics.git/internal/model"
)

type CustomMux struct {
	*http.ServeMux
}

func (m *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path == "/update/" {
	// 	updateHandler(w, r)
	// } else {
	// 	http.NotFound(w, r)
	// }
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(models.Storage); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Create custom mux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/update/", handler.UpdateHandler)
	mux.HandleFunc("/", mainHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
