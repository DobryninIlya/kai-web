package handler

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	handler "main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewScoreListHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getScoreList.NewScoreListHandler"
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка оценок: %v",
				path,
				err.Error(),
			)
			Respond(w, r, http.StatusNotFound, err)
			return
		}
		scoreElementList, err := handler.GetScoresStruct(scoreInfo.Faculty, scoreInfo.Course, scoreInfo.GroupId, scoreInfo.Idcard, scoreInfo.Studentid)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка структур БРС : %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}
		resultString, err := handler.GetScoreList(scoreElementList)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка оценок : %v",
				path,
				err.Error(),
			)
		}
		Respond(w, r, http.StatusOK, []byte(resultString))

	}
}
