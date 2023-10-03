package handler

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewLessonsHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getLesson.NewLessonHandler"
		uId := chi.URLParam(r, "uId")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		params := r.URL.Query()
		margin := params.Get("margin")
		marginI, err := strconv.Atoi(margin)
		if err != nil || uIdI < 0 {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение user: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}
		lessons, _, err := store.Schedule().GetCurrentDaySchedule(user.Group, marginI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение расписания на конкретный день: %v",
				path,
				err.Error(),
			)
		}
		lessonsDeleted, err := store.Schedule().GetDeletedLessonsByGroup(user.Group)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение расписания на конкретный день: %v",
				path,
				err.Error(),
			)
		}
		data := tools.GetLessonList(lessons, lessonsDeleted)
		Respond(w, r, http.StatusOK, []byte(data))
	}
}
