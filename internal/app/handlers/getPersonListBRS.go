package handler

import (
	"main/internal/app/tools"
	"net/http"
)

func NewPersonHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		pFac := params.Get("p_fac")
		pKurs := params.Get("p_kurs")
		pGroup := params.Get("p_group")
		if pFac == "" || pKurs == "" {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		result := tools.GetPersonListBRS(pFac, pKurs, pGroup)
		Respond(w, r, http.StatusOK, []byte(result))
	}
}
