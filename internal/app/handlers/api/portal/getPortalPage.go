package api_handler

import (
	"main/internal/app/tools"
	"net/http"
)

func NewPortalPageHandler(secretKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.portal.NewPortalPageHandler"
		paramsRaw := "?" + r.URL.RawQuery
		w.Write(tools.GetRegistrationPortal(paramsRaw))
	}
}
