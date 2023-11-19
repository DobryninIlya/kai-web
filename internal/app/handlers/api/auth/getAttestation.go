package auth

import (
	"context"
	"github.com/sirupsen/logrus"
	"main/internal/app/authorization"
	"main/internal/app/firebase"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
)

func NewAttestationHandler(ctx context.Context, store sqlstore.StoreInterface, log *logrus.Logger, fbAPI firebase.FirebaseAPIInterface, auth authorization.AuthorizationInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.auth.NewAttestationHandler"
		params := r.URL.Query()
		token := params.Get("token")
		client, err, _ := store.API().CheckToken(token)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		attestation, err := auth.GetAttestationList(token, client)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s: error while getting attestation list: %s",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}

		result := struct {
			Attestation []model.Discipline `json:"attestation"`
		}{
			Attestation: attestation,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
