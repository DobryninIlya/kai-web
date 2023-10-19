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

func NewHandleVKUpdateHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewHandleVKUpdateHandler"
		var upd model.VKUpdate
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка чтения body: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, h.ErrBadPayload)
			return
		}
		err = json.Unmarshal(body, &upd)
		if err != nil {
			log.Logf(
				logrus.WarnLevel,
				"%s : Ошибка анмарщалинга: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		if upd.Type == "confirmation" {
			result := store.API().AddAuthor(upd.Object.FromId)
			w.Write([]byte(store.API().ConfirmationCode))
			if result {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if upd.Type == "wall_post_new" {
			w.Write([]byte("ok"))
			if err := store.API().ParseNews(upd, log); err != nil {
				log.Logf(
					logrus.WarnLevel,
					"%v : Новость не сохранена %v",
					path,
					err,
				)
				return
			}
			return
		}
		w.Write([]byte("ok"))
	}
}
