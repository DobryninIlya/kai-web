package tasks

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewDeleteTaskHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeTask.NewDeleteTaskHandler"
		params := r.URL.Query()
		token := params.Get("token")
		idStr := chi.URLParam(r, "ID")
		groupnameStr := params.Get("groupname")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadParams)
			return
		}
		groupname, err := strconv.Atoi(groupnameStr)
		if err != nil || groupname < 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadParams)
			return
		}
		user, err, _ := store.API().CheckToken(token)
		if err != nil || user.Groupname != groupname {
			h.ErrorHandlerAPI(w, r, http.StatusForbidden, h.ErrForbidden)
			return
		}
		err = store.Task().Delete(id, user.Groupname)
		if err != nil {
			log.Error(path, err)
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, h.ErrInternal)
			return
		}
		h.RespondAPI(w, r, http.StatusOK, struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		})

	}
}
