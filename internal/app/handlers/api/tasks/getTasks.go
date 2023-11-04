package tasks

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewGetTaskHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeTask.NewTaskHandler"
		groupname := chi.URLParam(r, "groupname")
		groupnameInt, err := strconv.Atoi(groupname)
		if err != nil || groupnameInt < 0 || groupnameInt > 100000 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadParams)
			return
		}
		tasks, err := store.Task().GetAll(groupnameInt)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения заданий: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusOK, err)
			return
		}
		h.RespondAPI(w, r, http.StatusOK, struct {
			Tasks []model.Task `json:"tasks"`
		}{
			Tasks: tasks,
		})
	}
}
