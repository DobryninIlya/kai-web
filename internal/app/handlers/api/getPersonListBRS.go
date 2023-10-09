package api_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/tools"
	"net/http"
)

func NewPersonHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getPersonListBRS.NewPersonHandler"
		params := r.URL.Query()
		pFac := params.Get("p_fac")
		pKurs := params.Get("p_kurs")
		pGroup := params.Get("p_group")
		if pFac == "" || pKurs == "" {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		resultList, err := tools.GetPersonListBRS(pFac, pKurs, pGroup)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения групп БРС: %v",
				path,
				err.Error(),
			)
		}
		var result tools.GroupResultAnswer
		json.Unmarshal(resultList, &result)
		h.RespondAPI(w, r, http.StatusOK, struct {
			Groups map[string]string `json:"patronymics"`
		}{result.Result})
	}
}
