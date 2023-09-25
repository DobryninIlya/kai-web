package handler

import (
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewScoreHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		pFac, err := strconv.Atoi(params.Get("p_fac"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		pKurs, err := strconv.Atoi(params.Get("p_kurs"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		pGroup, err := strconv.Atoi(params.Get("p_group"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		pStud, _ := strconv.Atoi(params.Get("p_stud"))

		pZach, err := strconv.Atoi(params.Get("p_zach"))
		if pFac <= 0 || pKurs <= 0 || pGroup <= 0 || pStud <= 0 || pZach <= 0 || err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		result, err := tools.GetScores(pFac, pKurs, pGroup, pZach, pStud)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errUserNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		Respond(w, r, http.StatusOK, result)
	}
}
