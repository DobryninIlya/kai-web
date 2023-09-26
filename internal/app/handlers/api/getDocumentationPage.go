package api_handler

import (
	h "main/internal/app/handlers"
	"main/internal/app/tools"
	"net/http"
)

func NewDocumentationPageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := tools.GetDocumentationPage()
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
