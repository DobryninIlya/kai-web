package api_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
)

func NewMakeNewsHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewMakeNewsHandler"
		var res model.News
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение body: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrInternal)
			return
		}
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка анмаршалинга body: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		id, err := store.API().MakeNews(res)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка создания новости: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusOK, err)
		}
		h.RespondAPI(w, r, http.StatusOK, struct {
			Id int `json:"id"`
		}{
			Id: id,
		})
	}
}
