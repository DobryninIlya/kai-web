package handler

import (
	"encoding/json"
	"io"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewRegistrationHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res model.RegistrationData
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errBadPayload)
			return
		}
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errBadPayload)
			return
		}
		var groupId int
		var login string
		groupReal, err := strconv.Atoi(res.Identificator)
		if err == nil {
			login = ""
			groupId, _ = store.Schedule().GetIdByGroup(groupReal)
			if groupId == 0 {
				errorHandler(w, r, http.StatusBadRequest, errBadID)
				return
			}
		} else {
			login = res.Identificator
		}

		u := &model.User{
			ID:        res.VkId,
			Group:     groupId,
			GroupReal: groupReal,
			Role:      int8(res.Role),
			Login:     login,
		}
		//if val, err := service.MakeRegistration(res); val {
		if err := store.User().Create(u); err == nil {
			respond(w, r, http.StatusOK, []byte("{\"status\": \"ok\"}"))
			return
		} else {
			errorHandler(w, r, http.StatusBadRequest, errCantCreate)
		}

	}
}
