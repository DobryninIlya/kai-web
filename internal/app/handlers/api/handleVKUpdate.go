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
			store.API().AddAuthor(upd.GroupID)
			w.Write([]byte(store.API().ConfirmationCode))
			w.WriteHeader(http.StatusOK)
			return
		} else if upd.Type == "wall.post" {
			return
		}
	}
}
