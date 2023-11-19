package api_handler

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/authorization"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
)

func NewAttestationPageHandler(store sqlstore.StoreInterface, log *logrus.Logger, auth authorization.AuthorizationInterface, secretKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.portal.NewPortalPageHandler"
		url := r.URL.Query()
		tgID := url.Get("tg_id")
		uid := "tg" + tgID
		client, err := store.API().GetClient(uid)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		sign := GetSignForURLParams(r.URL.Query(), secretKey)
		list, err := auth.GetAttestationList(uid, client)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write(tools.GetAttestationPage(list, tgID, sign))
	}
}
