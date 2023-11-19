package api_handler

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"main/internal/app/authorization"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewAttestationElementPageHandler(store sqlstore.StoreInterface, log *logrus.Logger, auth authorization.AuthorizationInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.portal.NewPortalPageHandler"
		url := r.URL.Query()
		tgID := url.Get("tg_id")
		//att_element := url.Get("att_element")
		attElement := chi.URLParam(r, "id")
		attElementID, err := strconv.Atoi(attElement)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		uid := "tg" + tgID
		client, err := store.API().GetClient(uid)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s: error while getting client from db: %s",
				path, err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		list, err := auth.GetAttestationList(uid, client)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s: error while getting attestation list: %s",
				path, err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		result := list[attElementID]
		fmt.Println(result)
		w.Write(tools.GetAttestationElementPage(result))
		//h.RespondAPI(w, r, http.StatusOK, result)
	}
}
