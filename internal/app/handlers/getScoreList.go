package handler

import (
	"log"
	"main/internal/app/store/sqlstore"
	handler "main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewScoreListHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		if err != nil {
			log.Printf("Ошибка получения списка оценок, %v", err)
			Respond(w, r, http.StatusNotFound, err)
			return
		}
		scoreElementList, err := handler.GetScoresStruct(scoreInfo.Faculty, scoreInfo.Course, scoreInfo.GroupId, scoreInfo.Idcard, scoreInfo.Studentid)
		resultString := handler.GetScoreList(scoreElementList)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errUserNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		Respond(w, r, http.StatusOK, []byte(resultString))

	}
}
