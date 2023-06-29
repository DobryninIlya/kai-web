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
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		if err != nil || uIdI < 0 {
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errUserNotFound)
			return
		}
		lessons, _ := store.Schedule().GetCurrentDaySchedule(user.Group, marginI)
		data := tools.GetLessonList(lessons)
		respond(w, r, http.StatusOK, []byte(data))
	}
}
