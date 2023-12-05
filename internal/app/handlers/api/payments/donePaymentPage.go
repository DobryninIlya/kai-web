package api_handler

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"main/internal/payments"
	"net/http"
	"strings"
)

func NewDonePaymentPageHandler(log *logrus.Logger, store sqlstore.StoreInterface, pay payments.Yokassa) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.payments.NewDonePaymentPageHandler"
		uid := chi.URLParam(r, "payment_id")
		transaction, err := store.API().GetPaymentRequestByUID(uid)
		extID := strings.TrimSpace(transaction.ExtID)
		result, err := pay.CheckPaymentRequest(extID)
		if err != nil {
			log.Log(logrus.ErrorLevel, path+": "+err.Error())
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		if result.Status == "succeeded" {
			err = store.API().MakePremiumStatus(extID)
			if err != nil && err != sqlstore.ErrTransaсtionAlreadyEnded {
				log.Log(logrus.ErrorLevel, path+": "+err.Error()+" paymentID: "+extID)
				h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
				return
			}
		} else {
			page := tools.GetStatusPaymentPage(result.Status)
			w.WriteHeader(http.StatusOK)
			w.Write(page)
			return
		}
		page, err := tools.GetDonePaymentTemplate()
		if err != nil {
			log.Log(logrus.ErrorLevel, path+": "+err.Error())
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(page))
	}
}
