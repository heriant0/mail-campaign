package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/heriant0/mail-campaign/config"
	"github.com/heriant0/mail-campaign/helper"
)

var cfg config.Config

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config")
	}
	cfg = config

	http.HandleFunc("/send", enableCORS(send))

	appPort := fmt.Sprintf(":%s", cfg.AppPort)
	log.Println("server running on port", appPort)
	err = http.ListenAndServe(appPort, nil)
	if err != nil {
		panic(err)
	}
}

func send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request helper.MailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	byteRequest, err := json.Marshal(request)
	if err != nil {
		helper.ResponseError(w, http.StatusInternalServerError, "error when convert request to json data", err.Error())
		return
	}

	dataBuffer := bytes.NewBuffer(byteRequest)
	response, err := http.Post(cfg.BaseUrl, "application/json", dataBuffer)
	if err != nil {
		helper.ResponseError(w, http.StatusInternalServerError, "error when convert to buffer data", err.Error())
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("error saat parsing response body with error : %v", err.Error())
			return
		}
		log.Println("Response :", string(responseBytes))
		return
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("error saat parsing response body with error : %v", err.Error())
		return
	}
	log.Println("Response :", string(responseBytes))

	helper.ResponseSuccess(w, http.StatusOK, "email sent successfully")
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the actual handler
		handler(w, r)
	}
}
