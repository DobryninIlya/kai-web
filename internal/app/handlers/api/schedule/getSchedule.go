package schedule

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewScheduleHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getSchedule.NewScheduleHandler"
		groupId := chi.URLParam(r, "groupid")
		groupIdI, err := strconv.Atoi(groupId)
		if err != nil || groupId == "" {
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка получения параметров url запроса: %v",
					path,
					err.Error(),
				)
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		if err != nil || groupIdI <= 0 {
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка получения параметров url запроса: %v",
					path,
					err.Error(),
				)
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		lessons, err := store.Schedule().GetScheduleByGroup(groupIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения расписания группы: %v",
				path,
				err.Error(),
			)
		}
		resultLessons, err := store.Schedule().GetScheduleWithDeletedLessons(lessons, groupIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения расписания группы: %v",
				path,
				err.Error(),
			)
			return
		}
		result := struct {
			Lessons model.Schedule `json:"schedule"`
		}{
			resultLessons,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
