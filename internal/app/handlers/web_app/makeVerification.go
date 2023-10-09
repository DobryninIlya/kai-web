package web_app

import (
	"main/internal/app/tools"
	"net/http"
)

func NewVerificationTemplate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Respond(w, r, http.StatusOK, tools.GetRegistrationIDcard())
	}
}
