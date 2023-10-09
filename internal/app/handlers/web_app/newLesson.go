package web_app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewCreateLessonHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil || body == nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadPayload)
			return
		}
		var lesson model.LessonNew
		err = json.Unmarshal(body, &lesson)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadPayload)
			return
		}

		user, err := store.User().Find(uIdI)

		scoreInfo, err := store.Verification().GetPersonInfoScore(uIdI)
		if err != nil || scoreInfo.GroupId == 0 {
			log.Printf("Не удалось пометить занятие как удаленное: %v", err)
			ErrorHandler(w, r, http.StatusForbidden, errors.New(fmt.Sprintf("Не удалось пометить занятие как удаленное: %v", err)))
			return
		}
		_, err = store.Schedule().NewLesson(*user, lesson)
		if err != nil {
			log.Printf("Не удалось пометить занятие как удаленное: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Не удалось пометить занятие как удаленное: %v", err)))
			return
		}
		Respond(w, r, http.StatusOK, nil)
	}
}
