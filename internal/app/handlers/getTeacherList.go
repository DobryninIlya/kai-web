package handler

import (
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewTeacherHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errUserNotFound)

			return
		}
		teachers := store.Schedule().GetTeacherListStruct(user.Group)
		data := tools.GetTeacherList(teachers)
		respond(w, r, http.StatusOK, []byte(data))
	}
}
