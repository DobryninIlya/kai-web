package handler

import (
	"github.com/go-chi/chi"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewLessonsHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uId := chi.URLParam(r, "uId")
		uIdI, err := strconv.Atoi(uId)
		params := r.URL.Query()
		margin := params.Get("margin")
		marginI, err := strconv.Atoi(margin)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		if err != nil || uIdI < 0 {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}
		lessons, _ := store.Schedule().GetCurrentDaySchedule(user.Group, marginI)
		lessonsDeleted, _ := store.Schedule().GetDeletedLessonsByGroup(user.Group)
		data := tools.GetLessonList(lessons, lessonsDeleted)
		Respond(w, r, http.StatusOK, []byte(data))
	}
}
