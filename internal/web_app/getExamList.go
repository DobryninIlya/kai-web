package handler

import (
	"main/internal/database"
	"main/internal/handler"
	"net/http"
	"strconv"
)

func NewExamHandler(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		exam := service.GetExamListStruct(uIdI)
		data := handler.GetExamList(exam)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}
}
