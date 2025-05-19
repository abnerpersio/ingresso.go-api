package responses

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

func SendSuccess(w http.ResponseWriter, data ResponseData, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		SendError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(ResponseData{Message: message})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseData{Message: "Internal Server Error"})
		return
	}
}
