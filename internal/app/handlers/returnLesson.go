package handler

import (
	"errors"
	"fmt"
	"log"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewReturnLessonHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		user, err := store.User().Find(uIdI)

		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		if err != nil || scoreInfo.GroupId == 0 {
			log.Printf("Не удалось пометить занятие как удаленное: %v", err)
			ErrorHandler(w, r, http.StatusForbidden, errors.New(fmt.Sprintf("Не вернуть занятие в расписание: %v", err)))
			return
		}
		_, err = store.Schedule().ReturnDeletedLesson(*user, lessonIdI, uniqString)
		if err != nil {
			log.Printf("Не вернуть занятие в расписание: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Не вернуть занятие в расписание: %v", err)))
			return
		}
		Respond(w, r, http.StatusOK, nil)
	}
}
