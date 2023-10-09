package api_handler

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewTeachersHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getTeachers.NewTeachersHandler"
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
		teachers, err := store.Schedule().GetTeacherListStruct(groupIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка структур преподавателей: %v",
				path,
				err.Error(),
			)
		}
		result := struct {
			Teachers []model.Prepod `json:"teachers"`
		}{
			teachers,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
