package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
)

func NewNewsCreateFormHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getNewsCreateForm.NewNewsCreateFormHandler"
		page, err := tools.GetNewsCreatePage()
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка отрисовки страницы: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		h.Respond(w, r, http.StatusOK, page)
	}
}
