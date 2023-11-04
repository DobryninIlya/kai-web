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

func NewLessonsHandler(store sqlstore.StoreInterface, log *logrus.Logger, weekParity int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getLesson.NewLessonHandler"
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
		params := r.URL.Query()
		margin := params.Get("margin")
		marginI := 0
		if margin != "" {
			marginI, err = strconv.Atoi(margin)
		}
		if groupIdI <= 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		lessons, _, err := store.Schedule().GetCurrentDaySchedule(groupIdI, marginI, weekParity)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения текущего расписания на день : %v",
				path,
				err.Error(),
			)
		}
		result := struct {
			Schedule []model.Lesson `json:"schedule"`
		}{
			lessons,
		}
		h.RespondAPI(w, r, http.StatusOK, result)
	}
}
