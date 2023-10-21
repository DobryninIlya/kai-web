package api_handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/firebase"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strings"
)

func NewRegistrationHandler(ctx context.Context, store sqlstore.StoreInterface, log *logrus.Logger, fbAPI *firebase.FirebaseAPI) func(w http.ResponseWriter, r *http.Request) {
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
		if err != nil || res.DeviceTag == "" || res.UID == "" {
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
		if len(res.DeviceTag) > 16 || len(res.UID) > 35 {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка создания токена. Длина deviceTag или deviceId превышена : %v",
				path,
				h.ErrLongData.Error(),
			)
		}
		token, err := store.API().RegistrationToken(ctx, &res, fbAPI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения токена: %v",
				path,
				err.Error(),
			)
			if err == sqlstore.ErrUserNotFound {
				h.ErrorHandlerAPI(w, r, http.StatusNotFound, errors.New("пользователь не найден. Проверьте UID"))
				return
			}
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
