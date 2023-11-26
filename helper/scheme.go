package helper

import (
	"encoding/json"
	"net/http"
)

type MailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
	Type    string   `json:"type"`
}

type ResponseBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func ResponseError(w http.ResponseWriter, statusCode int, message string, errorMessage string) {
	response := ResponseBody{
		Success: false,
		Message: message,
		Error:   errorMessage,
	}

	writeJsonResponse(w, statusCode, response)
}

func ResponseSuccess(w http.ResponseWriter, statusCode int, message string) {
	response := ResponseBody{
		Success: true,
		Message: message,
	}

	writeJsonResponse(w, statusCode, response)
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, response ResponseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error endcode json response", http.StatusInternalServerError)
	}
}
