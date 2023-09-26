package api_handler

import (
	"github.com/go-chi/chi"
	h "main/internal/app/handlers"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewLessonsHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		groupId := chi.URLParam(r, "groupid")
		groupIdI, err := strconv.Atoi(groupId)
		if err != nil || groupId == "" {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		params := r.URL.Query()
		margin := params.Get("margin")
		marginI := 0
		if margin != "" {
			marginI, err = strconv.Atoi(margin)
		}
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		if err != nil || groupIdI <= 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		lessons, _ := store.Schedule().GetCurrentDaySchedule(groupIdI, marginI)
		h.RespondAPI(w, r, http.StatusOK, lessons)
	}
}
