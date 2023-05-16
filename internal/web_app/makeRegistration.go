package handler

import (
	"encoding/json"
	"io"
	"log"
	"main/internal/database"
	"net/http"
)

func NewRegistrationHandler(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res database.RegistrationData
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		if val, err := service.MakeRegistration(res); val {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("{\"status\": \"ok\"}"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		}

	}
}
