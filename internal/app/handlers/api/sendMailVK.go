package api_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/mailer"
	"main/internal/app/store/sqlstore"
	"main/internal/app/vk_api"
	"net/http"
)

func NewSendMailVKHandler(store sqlstore.StoreInterface, log *logrus.Logger, mail *mailer.Mailing) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewSendMailVKHandler"
		url := r.URL.Query()
		Type := url.Get("type")
		Label := url.Get("label")
		Link := url.Get("link")
		Payload := url.Get("payload")
		Message := url.Get("message")
		var keyboard vk_api.Keyboard
		if Type == "open_link" {
			if Label == "" || Link == "" || Message == "" {
				h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
				return
			}
			keyboard = vk_api.GetInlineLinkButtonVK(Type, Label, Link, Payload)
		} else {
			h.ErrorHandlerAPI(w, r, http.StatusNotImplemented, h.ErrNotImplemented)
			return
		}
		keyboardData, err := json.Marshal(keyboard)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка маршалинга: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		mail.GetVkClients()
		recipientCount := mail.SendMailVK(Message, string(keyboardData))

		h.RespondAPI(w, r, http.StatusOK, struct {
			RecipientCount int `json:"recipient_count"`
		}{recipientCount})
	}
}
