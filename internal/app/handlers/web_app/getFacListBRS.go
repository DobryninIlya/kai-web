package web_app

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/tools"
	"net/http"
)

func NewFacHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getFadcListBRS.NewFacHandler"
		result, err := tools.GetFacultiesListBRS()
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка факультетов для БРС: %v",
				path,
				err.Error(),
			)
		}
		Respond(w, r, http.StatusOK, result)
	}
}
