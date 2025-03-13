package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(writer http.ResponseWriter, code int, payload any) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON respone: %v", err)
		// HTTP status: 500
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(payloadJSON)
}

// 400
func RespondErrorWithJSON(writer http.ResponseWriter, code int, messageErr any) {
	if code >= http.StatusInternalServerError {
		// HTTP status: 500
		log.Printf("Responding with 5xx error: %v", messageErr)
	}

	RespondWithJSON(writer, code, messageErr)
}
