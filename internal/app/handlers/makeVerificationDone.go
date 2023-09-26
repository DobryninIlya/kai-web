package handler

import (
	"encoding/json"
	"io"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewVerificationDoneTemplate(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		idStr := params.Get("vk_user_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || body == nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		var ver model.VerificationParams
		err = json.Unmarshal(body, &ver)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadPayload)
			return
		}
		groupId, _ := store.Schedule().GetIdByGroup(ver.Groupname)
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
			ErrorHandler(w, r, http.StatusNotFound, err)
			return
		}

		err = store.User().MakeVerification(&ver, u)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, err)
			return
		}
		Respond(w, r, http.StatusCreated, nil)

	}
}
