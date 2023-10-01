package api_handler

import (
	"github.com/go-chi/chi"
	h "main/internal/app/handlers"
	"main/internal/app/tools"
	"net/http"
)

func NewDocumentationOtherPageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		if page == "" {
			page = "main"
		}
		result, err := tools.GetDocumentationPage(page)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return

	}
}
