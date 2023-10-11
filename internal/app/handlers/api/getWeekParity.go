package api_handler

import (
	h "main/internal/app/handlers/web_app"
	"net/http"
)

func NewWeekParityHandler(weekParity int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewWeekParityHandler"
		result := struct {
			WeekParity int `json:"week_parity"`
		}{
			weekParity,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
