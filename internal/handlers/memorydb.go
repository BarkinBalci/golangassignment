package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/BarkinBalci/golangassignment/pkg/memorydb"
)

var db = memorydb.NewInMemoryDB()

func MemoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreate(w, r)
	case http.MethodGet:
		handleFetch(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Printf("Error reading request body %v \n", err)
		return
	}

	var requestBody map[string]interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error unmarshaling request body", http.StatusBadRequest)
		log.Printf("Error unmarshaling request body %v\n", err)
		return
	}

	for key, value := range requestBody {
		db.Create(key, value)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Data created Successfully"))
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	data := db.FetchAll()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error encoding data as JSON", http.StatusInternalServerError)
		log.Printf("Error encoding data as JSON %v \n", err)
		return
	}
}
