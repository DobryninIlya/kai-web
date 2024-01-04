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

func NewDeleteLessonHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.deleteLesson.NewDeleteLessonHandler"
		params := r.URL.Query()
		token := params.Get("token")
		client, err, _ := store.API().CheckToken(token)
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
			h.ErrorHandler(w, r, http.StatusBadRequest, h.ErrBadParams)
			return
		}

		_, err = store.Schedule().MarkDeletedLesson(client.UID, client.Groupname, lessonIdI, uniqString, "mob")
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка попытки пометить задание как удаленное : %v",
				path,
				err.Error(),
			)
			h.ErrorHandler(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Не удалось пометить занятие как удаленное: %v", err)))
			return
		}
		h.Respond(w, r, http.StatusOK, nil)
	}
}
