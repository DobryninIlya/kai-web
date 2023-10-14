package api_handler

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewNewsHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewNewsHandler"
		newsIdS := chi.URLParam(r, "newsId")
		newsId, err := strconv.Atoi(newsIdS)
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
		news, err := store.API().GetNewsById(newsId)
		news.Id = newsId
		page, err := tools.GetNewsPage(news)
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
