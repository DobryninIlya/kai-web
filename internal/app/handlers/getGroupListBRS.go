package handler

import (
	"main/internal/app/tools"
	"net/http"
)

func NewGroupsHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		pFac := params.Get("p_fac")
		pKurs := params.Get("p_kurs")
		if pFac == "" || pKurs == "" {
			errorHandler(w, r, http.StatusBadRequest, errBadID)
			return
		}
		result := tools.GetGroupListBRS(pFac, pKurs)
		respond(w, r, http.StatusOK, result)
	}
}
