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

func NewPostProfilePhotoHandler(ctx context.Context, store sqlstore.StoreInterface, log *logrus.Logger, fbAPI firebase.FirebaseAPIInterface, auth authorization.AuthorizationInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.auth.NewAboutInfoHandler"
		params := r.URL.Query()
		token := params.Get("token")
		client, err, _ := store.API().CheckToken(token)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		file, _, err := r.FormFile("file")
		err = auth.ChangeProfilePhoto(client.UID, client, file)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		result := struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
