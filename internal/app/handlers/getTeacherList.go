package handler

import (
	model "main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

var teachersNull = make([]model.Prepod, 1)

func init() {
	teachersNull[0] = model.Prepod{
		LessonType: nil,
		Name:       " Нет данных",
		Lesson:     "",
	}
}

func NewTeacherHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errUserNotFound)

			return
		}
		teachers := store.Schedule().GetTeacherListStruct(user.Group)
		if len(teachers) == 0 {
			teachers = teachersNull
		}
		data := tools.GetTeacherList(teachers)
		respond(w, r, http.StatusOK, []byte(data))
	}
}
