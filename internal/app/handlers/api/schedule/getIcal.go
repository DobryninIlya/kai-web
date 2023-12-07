package schedule

import (
	"github.com/sirupsen/logrus"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/icalendar"
	"main/internal/app/store/sqlstore"
	"net/http"
	"strconv"
)

func NewIcalHandler(store sqlstore.StoreInterface, log *logrus.Logger, weekParity int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.getLesson.NewLessonHandler"
		params := r.URL.Query()
		userID := params.Get("user_id")
		userIDI, err := strconv.Atoi(userID)
		if err != nil || userID == "" {
			if err != nil {
				log.Logf(
					logrus.ErrorLevel,
					"%s : Ошибка получения параметров url запроса: %v",
					path,
					err.Error(),
				)
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		margin := params.Get("margin")
		marginI := 0
		if margin != "" {
			marginI, err = strconv.Atoi(margin)
		}
		if userIDI <= 0 {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		if !store.API().CheckPremiumByUID("tg" + userID) {
			if marginI > 3 {
				h.ErrorHandlerAPI(w, r, http.StatusForbidden, h.ErrForbidden)
				return
			}
		}

		user, err := store.API().GetTelegramUserByID(userIDI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения пользователя: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrBadID)
			return
		}
		icalFile, err := icalendar.GenerateICalendar(store.Schedule(), user.GroupID, marginI, weekParity, marginI)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка получения текущего расписания на день : %v",
				path,
				err.Error(),
			)
		}
		w.Header().Set("Content-Type", "text/calendar")
		w.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
		w.Write(icalFile)
	}
}
