package handler

import (
	"main/internal/database"
	"main/internal/handler"
	"net/http"
	"strconv"
)

func New(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		idStr := params.Get("vk_user_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if service.IsRegistredUser(id) {
			data := handler.GetMainView()
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		} else {
			data := handler.GetRegistrationView()
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}
}
