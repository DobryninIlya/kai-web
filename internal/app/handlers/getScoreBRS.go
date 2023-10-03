package handler

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewScoreHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getScoreBRS.NewScoreHandler"
		params := r.URL.Query()
		pFac, err := strconv.Atoi(params.Get("p_fac"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		pKurs, err := strconv.Atoi(params.Get("p_kurs"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		pGroup, err := strconv.Atoi(params.Get("p_group"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		pStud, _ := strconv.Atoi(params.Get("p_stud"))

		pZach, err := strconv.Atoi(params.Get("p_zach"))
		if pFac <= 0 || pKurs <= 0 || pGroup <= 0 || pStud <= 0 || pZach <= 0 || err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		result, err := tools.GetScores(pFac, pKurs, pGroup, pZach, pStud)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения БРС: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}
		Respond(w, r, http.StatusOK, result)
	}
}
