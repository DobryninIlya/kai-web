package handler

import (
	"main/internal/app/tools"
	"net/http"
)

func NewVerificationTemplate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, http.StatusOK, tools.GetRegistrationIDcard())
	}
}
