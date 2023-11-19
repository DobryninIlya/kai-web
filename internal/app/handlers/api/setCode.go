package api_handler

import (
	"errors"
	"github.com/go-chi/chi"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"net/http"
)

func NewSetCodeHandler(store *sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewSetCodeHandler"
		code := chi.URLParam(r, "code")
		if code == "" {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, errors.New("код подтверждения не указан"))
			return
		}
		(*store).API().SetConfirmationCode(code)
		h.RespondAPI(w, r, http.StatusOK, "токен успешно задан")
	}
}
