package api_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	h "main/internal/app/handlers"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strings"
)

func NewRegistrationHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewRegistrationHandler"
		var res model.ApiClient
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка чтения body: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		err = json.Unmarshal(body, &res)
		if err != nil || res.DeviceTag == "" || res.DeviceId == "" {
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка маршалинга body: %v",
					path,
					err.Error(),
				)
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		token, err := store.API().RegistrationToken(&res)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения токена: %v",
				path,
				err.Error(),
			)
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "ограничение уникальности") {
				h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrUniqueConstraint)
				return
			}
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		result := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
