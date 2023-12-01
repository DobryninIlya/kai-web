package api_handler

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/payments"
	"net/http"
)

func NewCheckPaymentRequestHandler(log *logrus.Logger, store sqlstore.StoreInterface, pay payments.Yokassa) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.payments.NewCheckPaymentRequestHandler"
		paymentID := chi.URLParam(r, "payment_id") // month, half-year, year
		result, err := pay.CheckPaymentRequest(paymentID)
		if err != nil {
			log.Log(logrus.ErrorLevel, path+": "+err.Error())
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		if result.Status == "succeeded" {
			err = store.API().MakePremiumStatus(paymentID)
			if err != nil && err != sqlstore.ErrTransationAlreadyEnded {
				log.Log(logrus.ErrorLevel, path+": "+err.Error()+" paymentID: "+paymentID)
				h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
				return
			}
		}
		h.RespondAPI(w, r, http.StatusOK, result.Status)
	}
}
