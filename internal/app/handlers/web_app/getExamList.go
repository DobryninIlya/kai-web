package web_app

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/store/parser"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tools"
	"net/http"
	"strconv"
)

func NewExamHandler(store sqlstore.StoreInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getExamList.NewExamHandler"
		params := r.URL.Query()
		uId := params.Get("vk_user_id")
		uIdI, err := strconv.Atoi(uId)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, ErrBadID)
			return
		}
		user, err := store.User().Find(uIdI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение расписания на конкретный день: %v",
				path,
				err.Error(),
			)
		}
		exam, err := parser.GetExamListStruct(user.Group)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получение списка экзаменов: %v",
				path,
				err.Error(),
			)
		}
		data := tools.GetExamList(exam)
		Respond(w, r, http.StatusOK, []byte(data))
	}
}
