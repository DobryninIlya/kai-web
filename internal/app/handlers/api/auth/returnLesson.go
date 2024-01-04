package auth

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewReturnLessonHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.returnLesson.NewReturnLessonHandler"
		params := r.URL.Query()
		token := params.Get("token")
		_, err, _ := store.API().CheckToken(token)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения токена: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		uniqString := params.Get("uniqstring")
		lessonId := params.Get("lesson_id")
		lessonIdI, err := strconv.Atoi(lessonId)
		if err != nil || lessonId == "" || uniqString == "" {
			h.ErrorHandler(w, r, http.StatusBadRequest, h.ErrBadPayload)
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
			h.ErrorHandler(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Не вернуть занятие в расписание: %v", err)))
			return
		}
		h.Respond(w, r, http.StatusOK, nil)
	}
}
