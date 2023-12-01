package api_handler

import (
	"main/internal/app/tools"
	"net/http"
)

func NewMakePaymentPageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.payments.NewDonePaymentPageHandler"
		page := tools.GetMakePaymentPage()
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	}
}
