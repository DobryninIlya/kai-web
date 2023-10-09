package web_app

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/tools"
	"net/http"
)

func NewGroupsHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getGroupListBRS.NewGroupsHandler"
		params := r.URL.Query()
		pFac := params.Get("p_fac")
		pKurs := params.Get("p_kurs")
		if pFac == "" || pKurs == "" {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		result, err := tools.GetGroupListBRS(pFac, pKurs)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения групп БРС: %v",
				path,
				err.Error(),
			)
		}
		Respond(w, r, http.StatusOK, result)
	}
}
