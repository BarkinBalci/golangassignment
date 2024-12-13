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

	var requestBody map[string]interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		log.Printf("Error unmarshaling request body %v\n", err)
		return
	}

	startDate, ok := requestBody["startDate"].(string)
	if !ok || startDate == "" {
		http.Error(w, "Missing or invalid start date", http.StatusBadRequest)
		log.Printf("Error missing or invalid start date %v\n", err)
		return
	}
	endDate, ok := requestBody["endDate"].(string)
	if !ok || endDate == "" {
		http.Error(w, "Missing or invalid end date", http.StatusBadRequest)
		log.Printf("Error missing or invalid end date %v\n", err)
		return
	}

	minCount, ok := requestBody["minCount"].(float64)
	if !ok {
		http.Error(w, "Missing or invalid minCount", http.StatusBadRequest)
		log.Printf("Error missing or invalid minCount %v\n", err)
		return
	}

	maxCount, ok := requestBody["maxCount"].(float64)
	if !ok {
		http.Error(w, "Missing or invalid maxCount", http.StatusBadRequest)
		log.Printf("Error missing or invalid maxCount %v\n", err)
		return
	}

	mongoClient, err := mongo.NewClient()
	if err != nil {
		http.Error(w, "Error connecting to MongoDB", http.StatusInternalServerError)
		log.Printf("Error connecting to MongoDB %v\n", err)
		return
	}
	records, err := mongoClient.FetchData(startDate, endDate, int(minCount), int(maxCount))

	if err != nil {
		if err.Error() == "invalid start date format" || err.Error() == "invalid end date format" {
			records = &mongo.FilteredRecords{
				Code:    1,
				Msg:     err.Error(),
				Records: []mongo.Record{},
			}
		} else {
			http.Error(w, "Error fetching data from MongoDB", http.StatusInternalServerError)
			log.Printf("Error fetching data from MongoDB %v \n", err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(records)
	if err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
}
