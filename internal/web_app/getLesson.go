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
		uId := chi.URLParam(r, "uId")
		uIdI, err := strconv.Atoi(uId)
		params := r.URL.Query()
		margin := params.Get("margin")
		marginI, err := strconv.Atoi(margin)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil || uIdI < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		lessons, _ := service.GetCurrentDaySchedule(uIdI, marginI)
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
