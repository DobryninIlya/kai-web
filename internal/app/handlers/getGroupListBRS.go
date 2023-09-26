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
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		result := tools.GetGroupListBRS(pFac, pKurs)
		Respond(w, r, http.StatusOK, result)
	}
}
