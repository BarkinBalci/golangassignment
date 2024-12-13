package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/BarkinBalci/golangassignment/pkg/mongo"
)

func MongoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Printf("Error reading request body %v\n", err)
		return
	}

	var requestBody map[string]string
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		log.Printf("Error unmarshaling request body %v\n", err)
		return
	}
	collectionName, ok := requestBody["collection_name"]
	if !ok || collectionName == "" {
		http.Error(w, "Missing collection name", http.StatusBadRequest)
		log.Printf("Error missing collection_name %v\n", err)
		return
	}

	mongoClient, err := mongo.NewMongoClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		log.Printf("Error creating mongo client, %v \n", err)
		return
	}

	results, err := mongoClient.FetchData(collectionName)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		log.Printf("Error Fetching data %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Error marshaling data", http.StatusInternalServerError)
		return
	}
}
