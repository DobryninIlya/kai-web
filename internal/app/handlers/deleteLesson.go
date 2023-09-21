package handler

import (
	"errors"
	"fmt"
	"log"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewDeleteLessonHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		lessonId := params.Get("lesson_id")
		lessonIdI, err := strconv.Atoi(lessonId)
		if err != nil || lessonId == "" {
			errorHandler(w, r, http.StatusBadRequest, errBadPayload)
			return
		}
		user, err := store.User().Find(uIdI)

		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		if err != nil || scoreInfo.GroupId == 0 {
			log.Printf("Не удалось пометить занятие как удаленное: %v", err)
			errorHandler(w, r, http.StatusForbidden, errors.New(fmt.Sprintf("Не удалось пометить занятие как удаленное: %v", err)))
			return
		}
		_, err = store.Schedule().MarkDeletedLesson(*user, lessonIdI)
		if err != nil {
			log.Printf("Не удалось пометить занятие как удаленное: %v", err)
			errorHandler(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Не удалось пометить занятие как удаленное: %v", err)))
			return
		}
		respond(w, r, http.StatusOK, nil)
	}
}
