package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/heriant0/mail-campaign/config"
	"github.com/heriant0/mail-campaign/helper"
	"gopkg.in/gomail.v2"
)

var cfg config.Config

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config")
	}

	cfg = config
	http.HandleFunc("/send", enableCORS(sendMail))

	mailPort := fmt.Sprintf(":%s", cfg.MailPort)
	log.Println("server running on port", mailPort)
	err = http.ListenAndServe(mailPort, nil)
	if err != nil {
		panic(err)
	}

}

func sendMail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		helper.ResponseError(w, http.StatusMethodNotAllowed, "method not allowd", "")
		return
	}

	var request helper.MailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	err := sendWithGomail(request)
	if err != nil {
		panic(err)
	}

	helper.ResponseSuccess(w, http.StatusOK, "email sent successfully")
}

func sendWithGomail(mailRequest helper.MailRequest) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", cfg.ConfigSenderName)
	mailer.SetHeader("To", mailRequest.To...)
	mailer.SetHeader("Subject", mailRequest.Subject)
	mailer.SetBody("text/html", mailRequest.Message)

	dialer := gomail.NewDialer(
		cfg.ConfigSmtpHost,
		cfg.ConfigSmtpPort,
		cfg.ConfigAuthEmail,
		cfg.ConfigAuthPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatalf("error when try to dial and send message : %v", err.Error())
		return err
	}

	return nil
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the actual handler
		handler(w, r)
	}
}
