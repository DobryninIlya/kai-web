package auth

import (
	"context"
	"github.com/sirupsen/logrus"
	"main/internal/app/authorization"
	"main/internal/app/firebase"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"net/http"
)

func NewProfilePhotoURLHandler(ctx context.Context, store sqlstore.StoreInterface, log *logrus.Logger, fbAPI firebase.FirebaseAPIInterface, auth authorization.AuthorizationInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.auth.NewAboutInfoHandler"
		params := r.URL.Query()
		token := params.Get("token")
		client, err, _ := store.API().CheckToken(token)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		userInfo, err := auth.GetProfilePhotoURL(token, client)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		result := struct {
			UserInfo string `json:"photo_url"`
		}{
			UserInfo: userInfo,
		}
		w.Header().Set("Content-Type", "application/json")
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
