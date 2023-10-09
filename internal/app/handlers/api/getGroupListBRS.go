package api_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"net/http"
)

import (
	"main/internal/app/tools"
)

func NewGroupsHandler(log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getGroupListBRS.NewGroupsHandler"
		params := r.URL.Query()
		pFac := params.Get("p_fac")
		pKurs := params.Get("p_kurs")
		if pFac == "" || pKurs == "" {
			log.Logf(
				logrus.WarnLevel,
				"%s : Ошибка получения параметров запроса: %v, %v",
				path,
				pFac,
				pKurs,
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		resultList, err := tools.GetGroupListBRS(pFac, pKurs)
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
			Groups map[string]string `json:"groups"`
		}{result.Result})
	}
}
