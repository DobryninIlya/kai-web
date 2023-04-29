package handler

import (
	"main/internal/database"
	"main/internal/handler"
	"net/http"
)

func New(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lessons := service.GetCurrentDaySchedule(23546, 3)
		data := handler.GetMainView(lessons)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
