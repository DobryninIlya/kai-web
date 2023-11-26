package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/parser"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewExamPageHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.portal.NewExamPageHandler"
		url := r.URL.Query()
		groupID := url.Get("groupid")
		group, err := strconv.Atoi(groupID)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		exam, err := parser.GetExamListStruct(group)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение списка экзаменов: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusNotFound, h.ErrBadID)
			return
		}
		w.Write(tools.GetExamPage(exam))
	}
}
