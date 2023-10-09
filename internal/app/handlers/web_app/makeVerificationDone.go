package web_app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewVerificationDoneTemplate(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.makeVerificationDone.NewVerificationDoneTemplate"
		params := r.URL.Query()
		idStr := params.Get("vk_user_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || body == nil {
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка получение user: %v",
					path,
					err.Error(),
				)
			}
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		var ver model.VerificationParams
		err = json.Unmarshal(body, &ver)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка анмаршалинга body: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadPayload)
			return
		}
		groupId, err := store.Schedule().GetIdByGroup(ver.Groupname)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение groupId: %v",
				path,
				err.Error(),
			)
		}
		if groupId == 0 {
			ErrorHandler(w, r, http.StatusNotFound, ErrUserNotFound)
			return
		}
		u := &model.User{
			ID:        id,
			Group:     groupId,
			GroupReal: ver.Group,
			Role:      int8(1),
		}

		_, err = tools.GetScores(ver.Faculty, ver.Course, ver.Group, ver.ID, ver.Student)
		if err != nil { // Если данные БРС получены, можно сохранять в базе
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение БРС: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusNotFound, err)
			return
		}

		err = store.User().MakeVerification(&ver, u)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка верификации user: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, err)
			return
		}
		Respond(w, r, http.StatusCreated, nil)

	}
}
