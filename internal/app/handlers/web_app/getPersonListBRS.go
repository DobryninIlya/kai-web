package web_app

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/tools"
	"net/http"
)

func NewPersonHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getPersonListBRS.NewPersonHandler"
		params := r.URL.Query()
		pFac := params.Get("p_fac")
		pKurs := params.Get("p_kurs")
		pGroup := params.Get("p_group")
		if pFac == "" || pKurs == "" {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		result, err := tools.GetPersonListBRS(pFac, pKurs, pGroup)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения ФИО БРС: %v",
				path,
				err.Error(),
			)
		}
		Respond(w, r, http.StatusOK, result)
	}
}
