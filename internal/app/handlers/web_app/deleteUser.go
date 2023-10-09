package web_app

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewDeleteUserHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.deleteUser.NewDeleteUserHandler"
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		err = store.User().Delete(uIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение user: %v",
				path,
				err.Error(),
			)
		}
		Respond(w, r, http.StatusOK, nil)
	}
}
