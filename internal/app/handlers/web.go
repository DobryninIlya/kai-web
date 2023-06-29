package handler

import (
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func New(store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			return
		} else {
			data := tools.GetRegistrationView()
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}
}
