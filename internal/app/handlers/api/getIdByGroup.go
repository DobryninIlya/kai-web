package api_handler

import (
	"github.com/go-chi/chi"
	h "main/internal/app/handlers"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

type answer struct {
	GroupId int `json:"group_id"`
}

func NewIdByGroupHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group := chi.URLParam(r, "group")
		groupI, err := strconv.Atoi(group)
		if err != nil || group == "" {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		if err != nil || groupI <= 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		groupID, error := store.Schedule().GetIdByGroup(groupI)
		if error != nil {
			h.ErrorHandlerAPI(w, r, http.StatusNotImplemented, h.ErrInternal)
			return
		}
		result := answer{GroupId: groupID}
		//data, _ := json.Marshal(result)
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
