package handler

import (
	"main/internal/app/tools"
	"net/http"
)

func NewFacHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result := tools.GetFacultiesListBRS()
		Respond(w, r, http.StatusOK, result)
	}
}
