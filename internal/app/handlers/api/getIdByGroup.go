package api_handler

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

type answer struct {
	GroupId int `json:"group_id"`
}

func NewIdByGroupHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getIdByGroup.NewIdByGroupHandler"
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
		groupID, err := store.Schedule().GetIdByGroup(groupI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения расписания: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusNotFound, h.ErrRecordNotFound)
			return
		}
		result := answer{GroupId: groupID}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
