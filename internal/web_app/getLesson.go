package handler

import (
	"github.com/go-chi/chi"
	"main/internal/database"
	"main/internal/handler"
	"net/http"
	"strconv"
)

type result struct {
	data string `json:"data"`
}

func NewLessonsHandler(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		offset := chi.URLParam(r, "offset")
		offsetI, err := strconv.Atoi(offset)
		if err != nil || offsetI < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		lessons, _ := service.GetCurrentDaySchedule(23546, offsetI)
		data := handler.GetLessonList(lessons)
		var d result
		d.data = data
		//value, _ := json.Marshal(&d)
		//w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}
}
