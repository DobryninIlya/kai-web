package api_handler

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/tools"
	"net/http"
)

func NewFacHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getFacListBRS.NewFacHandler"
		facList := tools.Faculties
		h.RespondAPI(w, r, http.StatusOK, struct {
			Faculties map[int]string `json:"faculties"`
		}{facList})
	}
}
