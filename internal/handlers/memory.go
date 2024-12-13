package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/BarkinBalci/golangassignment/pkg/memory"
)

var mem = memory.NewItem()

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

	var requestBody struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error unmarshaling request body", http.StatusBadRequest)
		log.Printf("Error unmarshaling request body %v\n", err)
		return
	}

	if requestBody.Key == "" || requestBody.Value == "" {
		http.Error(w, "Key or Value cannot be empty", http.StatusBadRequest)
		log.Printf("Key or Value cannot be empty %v\n", err)
		return
	}

	mem.Create(requestBody.Key, requestBody.Value)

	responseBody, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Error encoding the response body", http.StatusInternalServerError)
		log.Printf("Error encoding the response body %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Error parsing query parameters", http.StatusBadRequest)
		log.Printf("Error parsing query parameters %v \n", err)
		return
	}
	key := query.Get("key")

	if key == "" {
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		log.Println("Missing 'key' query parameter")
		return
	}

	value, ok := mem.Fetch(key)

	if !ok {
		http.Error(w, "Key not found in the database", http.StatusNotFound)
		log.Println("Key not found in the database")
		return
	}

	response := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   key,
		Value: value,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding data as JSON", http.StatusInternalServerError)
		log.Printf("Error encoding data as JSON %v \n", err)
		return
	}
}
