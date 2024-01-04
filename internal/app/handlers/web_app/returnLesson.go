package web_app

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewReturnLessonHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.returnLesson.NewReturnLessonHandler"
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		uniqString := params.Get("uniqstring")
		lessonId := params.Get("lesson_id")
		lessonIdI, err := strconv.Atoi(lessonId)
		if err != nil || lessonId == "" || uniqString == "" {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadPayload)
			return
		}
		_, err = store.User().Find(uIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения user: %v",
				path,
				err.Error(),
			)
		}
		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		if err != nil || scoreInfo.GroupId == 0 {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения верификационных данных : %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusForbidden, errors.New(fmt.Sprintf("Не вернуть занятие в расписание: %v", err)))
			return
		}
		_, err = store.Schedule().ReturnDeletedLesson(lessonIdI, uniqString)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка возврата занятия в расписание : %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Не вернуть занятие в расписание: %v", err)))
			return
		}
		Respond(w, r, http.StatusOK, nil)
	}
}
