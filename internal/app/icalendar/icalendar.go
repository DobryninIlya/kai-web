package icalendar

import (
	ical "github.com/arran4/golang-ical"
	"main/internal/app/store/sqlstore"
	"strings"
	"time"
)

func GenerateICalendar(scheduleRepository *sqlstore.ScheduleRepository, groupId int, margin int, weekParity int, numDays int) ([]byte, error) {
	cal := ical.NewCalendar()

	for i := 0; i < numDays; i++ {
		day := time.Now().AddDate(0, 0, margin+i)
		lessons, _, err := scheduleRepository.GetCurrentDaySchedule(groupId, margin+i, weekParity)
		if err != nil {
			return nil, err
		}

		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
		//dayEnd := dayStart.Add(24 * time.Hour)

		for _, lesson := range lessons {
			startTime, err := time.Parse("15:04", strings.TrimSpace(lesson.DayTime))
			if err != nil {
				return nil, err
			}
			duration, _ := time.ParseDuration("1h30m")

			event := cal.AddEvent(lesson.DisciplName)
			startAt := dayStart.Add(time.Duration(startTime.Hour()) * time.Hour).Add(time.Duration(startTime.Minute()) * time.Minute)
			endAt := startAt.Add(duration)
			event.SetStartAt(startAt)
			event.SetEndAt(endAt)
			event.SetLocation(lesson.BuildNum + " " + lesson.AudNum)
		}
	}

	return []byte(cal.Serialize()), nil
}
