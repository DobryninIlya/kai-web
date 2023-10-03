package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewRegistrationHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.makeRegistration.NewRegistrationHandler"
		var res model.RegistrationData
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение user: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadPayload)
			return
		}
		err = json.Unmarshal(body, &res)
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
		var groupId int
		var login string
		groupReal, err := strconv.Atoi(res.Identificator)
		if err == nil { // В таком случае ожидаем числовой ID
			login = ""
			groupId, err = store.Schedule().GetIdByGroup(groupReal)
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка получения расписания : %v",
					path,
					err.Error(),
				)
				ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
				return
			}
			if groupId == 0 {
				ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
				return
			}
		} else if groupReal == 0 {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		} else { // В таком случае ожидаем стринговый айди
			login = res.Identificator
		}

		u := &model.User{
			ID:        res.VkId,
			Group:     groupId,
			GroupReal: groupReal,
			Role:      int8(res.Role) + 1,
			Login:     login,
			Name:      "|имя не задано|",
		}
		//if val, err := service.MakeRegistration(res); val {
		if err := store.User().Create(u); err == nil {
			Respond(w, r, http.StatusOK, []byte("{\"status\": \"ok\"}"))
			return
		} else {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка создания user: %v",
				path,
				err.Error(),
			)
			ErrorHandler(w, r, http.StatusBadRequest, ErrCantCreated)
		}

	}
}
