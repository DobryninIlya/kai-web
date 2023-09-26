package api_handler

import (
	"github.com/go-chi/chi"
	h "main/internal/app/handlers"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewTeachersHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		groupId := chi.URLParam(r, "groupid")
		groupIdI, err := strconv.Atoi(groupId)
		if err != nil || groupId == "" {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		if err != nil || groupIdI <= 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		teachers := store.Schedule().GetTeacherListStruct(groupIdI)
		h.RespondAPI(w, r, http.StatusOK, teachers)
	}
}
