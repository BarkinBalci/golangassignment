package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/BarkinBalci/golangassignment/pkg/mongo"
)

type APIResponse struct {
	Code    int            `json:"code"`
	Msg     string         `json:"msg"`
	Records []mongo.Record `json:"records,omitempty"`
}

func MongoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, &APIResponse{Code: 1, Msg: "Method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: 2, Msg: "Error reading request body"})
		return
	}

	var requestBody map[string]interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: 3, Msg: "Error unmarshalling request body"})
		return
	}

	startDate, ok := requestBody["startDate"].(string)
	if !ok || startDate == "" {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: 4, Msg: "Missing or invalid start date"})
		return
	}

	endDate, ok := requestBody["endDate"].(string)
	if !ok || endDate == "" {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: 5, Msg: "Missing or invalid end date"})
		return
	}

	minCount, ok := requestBody["minCount"].(float64)
	if !ok {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: 6, Msg: "Missing or invalid minCount"})
		return
	}

	maxCount, ok := requestBody["maxCount"].(float64)
	if !ok {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: 7, Msg: "Missing or invalid maxCount"})
		return
	}

	mongoClient, err := mongo.NewClient()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &APIResponse{Code: 8, Msg: "Error connecting to MongoDB"})
		log.Printf("Error connecting to MongoDB: %v\n", err)
		return
	}

	records, err := mongoClient.FetchData(startDate, endDate, int(minCount), int(maxCount))
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &APIResponse{Code: -1, Msg: err.Error()})
		return
	}

	writeResponse(w, http.StatusOK, &APIResponse{Code: 0, Msg: "Success", Records: records})
}

func writeResponse(w http.ResponseWriter, statusCode int, response *APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
	}
}
