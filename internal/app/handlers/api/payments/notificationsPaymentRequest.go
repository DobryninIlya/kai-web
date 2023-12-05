package api_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/payments"
	"net/http"
)

func NewNotificationsPaymentRequestHandler(log *logrus.Logger, store sqlstore.StoreInterface, pay payments.Yokassa) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.payments.NewCheckPaymentRequestHandler"
		var notification payments.Notification
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Log(logrus.ErrorLevel, path+": "+err.Error())
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		json.Unmarshal(body, &notification)
		if notification.Event != "payment.succeeded" {
			h.RespondAPI(w, r, http.StatusOK, "not payment.succeeded")
			return
		}
		result, err := pay.CheckPaymentRequest(notification.Object.Id)
		if err != nil {
			log.Log(logrus.ErrorLevel, path+": "+err.Error())
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		if result.Status == "succeeded" {
			err = store.API().MakePremiumStatus(notification.Object.Id)
			if err != nil && err != sqlstore.ErrTransaсtionAlreadyEnded {
				log.Log(logrus.ErrorLevel, path+": "+err.Error()+" paymentID: "+notification.Object.Id)
				h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
				return
			}
		}
		h.RespondAPI(w, r, http.StatusOK, result.Status)
	}
}
