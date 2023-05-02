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
		uId := params.Get("vk_user_id")
		userId, err := strconv.Atoi(uId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		lessons, date := service.GetCurrentDaySchedule(userId, 0)
		//age := params.Get("age")
		data := handler.GetMainView(lessons, date)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
