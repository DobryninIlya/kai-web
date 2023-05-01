package database

import (
	"github.com/jackc/pgx"
	"log"
	"time"
)

type Service struct {
	conn *pgx.Conn
}

func (s *Service) GetScheduleByGroup(id int) Schedule {
	rows, err := s.conn.Query("SELECT shedule FROM saved_timetable WHERE groupp = $1", id)
	if err != nil {
		log.Println("Ошибка запроса расписания")
		log.Println(err)
	}
	var schedule string
	defer rows.Close()
	rows.Next()
	rows.Scan(&schedule)
	scheduleStruct := GetScheduleStruct([]byte(schedule))
	return scheduleStruct
}

func (s *Service) GetCurrentDaySchedule(group int, margin int) ([]Lesson, time.Time) {
	day := time.Now().AddDate(0, 0, margin)
	dayNum := day.Weekday()
	groupSchedule := s.GetScheduleByGroup(group)
	lessons := make([]Lesson, 2)
	switch {
	case dayNum == 1:
		lessons = groupSchedule.Day1
	case dayNum == 2:
		lessons = groupSchedule.Day2
	case dayNum == 3:
		lessons = groupSchedule.Day3
	case dayNum == 4:
		lessons = groupSchedule.Day4
	case dayNum == 5:
		lessons = groupSchedule.Day5
	case dayNum == 6:
		lessons = groupSchedule.Day6
	case dayNum == 7:
		lessons = []Lesson{}
	}

	return lessons, day
}

func NewService(pgConfig pgx.ConnConfig) *Service {
	conn, err := pgx.Connect(pgConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return &Service{
		conn: conn,
	}
}
