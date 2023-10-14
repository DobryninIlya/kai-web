package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewNewsPreviewsHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewNewsPreviewsHandler"
		url := r.URL.Query()
		countS := url.Get("count")
		count, err := strconv.Atoi(countS)
		if err != nil {
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
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrIncorrectParams)
		}
		news, err := store.API().GetNewsPreviews(count, offset)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка новостей превью: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		h.RespondAPI(w, r, http.StatusOK, struct {
			News []model.News `json:"news"`
		}{
			News: news,
		})
	}
}
