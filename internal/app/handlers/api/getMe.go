package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
)

func NewWhoIAmHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getMe.NewWhoIAmHandler"
		url := r.URL.Query()
		token := url.Get("token")
		client, err, _ := store.API().CheckToken(token)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения токена: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		client.EncryptedPassword = nil
		result := struct {
			model.ApiRegistration
		}{
			client,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
