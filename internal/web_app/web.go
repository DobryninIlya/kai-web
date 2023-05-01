package handler

import (
	"main/internal/database"
	"main/internal/handler"
	"net/http"
)

func New(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lessons, date := service.GetCurrentDaySchedule(23546, 3)
		//params := r.URL.Query()
		//name := params.Get("name")
		//age := params.Get("age")
		data := handler.GetMainView(lessons, date)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
