package handler

import (
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
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		scoreElementList, err := handler.GetScoresStruct(scoreInfo.Faculty, scoreInfo.Course, scoreInfo.GroupId, scoreInfo.Idcard, scoreInfo.Studentid)
		resultString := handler.GetScoreList(scoreElementList)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errUserNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		respond(w, r, http.StatusOK, []byte(resultString))

	}
}