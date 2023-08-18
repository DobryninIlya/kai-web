package handler

import (
	"log"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewDeleteUserHandler(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		err = store.User().Delete(uIdI)
		if err != nil {
			log.Printf("Не удалось удалить запись пользователя: %v", err)
		}
		respond(w, r, http.StatusOK, nil)
	}
}
