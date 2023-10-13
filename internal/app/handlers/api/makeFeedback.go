package api_handler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tg_api"
	"net/http"
)

const (
	peer            = -1001907805430
	messageTemplate = `*Сообщение*:
*device_id*=%v
*device_tage*=%v
*ver.* #%v
*text:*
_%v_
`
)

func NewFeedbackHandler(store sqlstore.StoreInterface, log *logrus.Logger, api *tg_api.APItg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeFeedback.NewFeedbackHandler"
		var res model.FeedbackClient
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
		err = json.Unmarshal(body, &res)
		if res.Version == "" || res.Text == "" {
			log.Logf(
				logrus.WarnLevel,
				"%s : Данные body пустые или некорректны",
				path,
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		url := r.URL.Query()
		tokenInfo, err := store.API().GetTokenInfo(url.Get("token"))
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения информации о токене: %v",
				path,
				err.Error(),
			)
		}
		message := fmt.Sprintf(messageTemplate, tokenInfo.UID, tokenInfo.DeviceTag, res.Version, res.Text)
		sendResult := api.SendMessageTG(log, peer, message, "", 592)
		if !sendResult {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка отправления обратной связи в телеграм-чат: %v, $v",
				path,
				err.Error(),
				sendResult,
			)
		}
		result := struct {
			Success bool `json:"success"`
		}{
			sendResult,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
