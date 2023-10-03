package handler

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func New(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.web.New"
		params := r.URL.Query()
		idStr := params.Get("vk_user_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		if user, err := store.User().Find(id); err == nil && user != nil {
			data := tools.GetMainView()
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение user: %v",
				path,
				err,
			)
			return
		} else {
			if err != nil {
				log.Printf("Error when user create: %v", err)
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка чтения body: %v",
					path,
					err,
				)
			}
			data := tools.GetRegistrationView()
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}
}
