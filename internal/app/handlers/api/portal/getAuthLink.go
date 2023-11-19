package api_handler

import (
	h "main/internal/app/handlers/web_app"
	"net/http"
)

func NewAuthLinkHandler(secretKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.portal.NewAuthLinkHandler"
		url := r.URL.Query()
		sign := GetSignForURLParams(url, secretKey)
		result := struct {
			Sign string `json:"sign"`
		}{
			Sign: sign,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
