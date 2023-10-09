package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewScoreHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getScoreBRS.NewScoreHandler"
		params := r.URL.Query()
		pFac, err := strconv.Atoi(params.Get("p_fac"))
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения параметров url запроса: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return

		}
		pKurs, err := strconv.Atoi(params.Get("p_kurs"))
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения параметров url запроса: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		pGroup, err := strconv.Atoi(params.Get("p_group"))
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения параметров url запроса: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		pStud, err := strconv.Atoi(params.Get("p_stud"))
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения параметров url запроса: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		pZach, err := strconv.Atoi(params.Get("p_zach"))
		if pFac <= 0 || pKurs <= 0 || pGroup <= 0 || pStud <= 0 || pZach <= 0 || err != nil {
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка получения параметров url запроса: %v",
					path,
					err.Error(),
				)
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		scoresList, err := tools.GetScoresStruct(pFac, pKurs, pGroup, pZach, pStud)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения БРС: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrUserNotFound)
			return
		}
		result := struct {
			Scores []tools.ScoreElement `json:"scores"`
		}{
			scoresList,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
