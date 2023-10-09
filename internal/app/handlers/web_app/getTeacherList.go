package web_app

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

//var teachersNull = make([]model.Prepod, 1)
//
//func init() {
//	teachersNull[0] = model.Prepod{
//		LessonType: nil,
//		Name:       " Нет данных",
//		Lesson:     "",
//	}
//}

func NewTeacherHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getTeachersList.NewTeacherHandler"
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения расписания на конкретный день: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}
		teachers, err := store.Schedule().GetTeacherListStruct(user.Group)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка преподавателей: %v",
				path,
				err.Error(),
			)
		}
		data, err := tools.GetTeacherList(teachers)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения списка преподавателей : %v",
				path,
				err.Error(),
			)
		}
		Respond(w, r, http.StatusOK, []byte(data))
	}
}
