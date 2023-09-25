package handler

import (
	"main/internal/app/store/parser"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewExamHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {

		}
		exam := parser.GetExamListStruct(user.Group)
		data := tools.GetExamList(exam)
		Respond(w, r, http.StatusOK, []byte(data))
	}
}
