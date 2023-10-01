package api_handler

import (
	"encoding/json"
	"io"
	h "main/internal/app/handlers"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strings"
)

func NewRegistrationHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res model.ApiClient
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		err = json.Unmarshal(body, &res)
		if err != nil || res.DeviceTag == "" || res.DeviceId == "" {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadPayload)
			return
		}
		token, err := store.API().RegistrationToken(&res)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "ограничение уникальности") {
				h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrUniqueConstraint)
				return
			}
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		result := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
