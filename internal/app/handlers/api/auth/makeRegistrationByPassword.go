package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/authorization"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strings"
)

func NewRegistrationByPasswordHandler(ctx context.Context, store sqlstore.StoreInterface, log *logrus.Logger, auth authorization.AuthorizationInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.auth.newRegistrationByPasswordHandler"
		var res model.ApiRegistration
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
		if err != nil || res.Login == "" || res.Password == "" {
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
		if len(res.UID) > 35 {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка создания токена. Длина deviceTag или deviceId превышена : %v",
				path,
				h.ErrLongData.Error(),
			)
		}
		token, err := store.API().RegistrationUserByPassword(ctx, &res, auth, res.Login, res.Password)
		if err != nil && token == "" {
			if strings.Contains(err.Error(), "UNIQUE constraint") || strings.Contains(err.Error(), "ограничение уникальности") ||
				strings.Contains(err.Error(), "unique constraint") || errors.Is(sql.ErrNoRows, err) {
				//h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrUniqueConstraint)
				token, err = store.API().GetTokenByUID(store.API().GenerateUID(res.Login, res.Password))
				result := struct {
					Token string `json:"token"`
				}{
					Token: token,
				}
				h.RespondAPI(w, r, http.StatusOK, result)
				return
			}
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения токена: %v",
				path,
				err.Error(),
			)
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
