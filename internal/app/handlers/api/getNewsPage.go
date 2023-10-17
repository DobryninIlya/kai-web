package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewNewsPageHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.geNewsPage.NewNewsPageHandler"
		url := r.URL.Query()
		countS := url.Get("count")
		count, err := strconv.Atoi(countS)
		if countS == "" {
			count = 20
		}
		if err != nil && countS != "" {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения newsId: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		offsetS := url.Get("offset")
		offset, err := strconv.Atoi(offsetS)
		if offsetS == "" {
			offset = 0
		}
		if err != nil && offsetS != "" {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения newsId: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		if count == 0 {
			count = 20
		}
		news, err := store.API().GetNewsPreviews(count, offset)
		page, err := tools.GetNewsPreviewsPage(news)
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
