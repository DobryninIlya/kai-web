package model

import (
	"github.com/jackc/pgx/pgtype"
	"time"
)

type ScheduleSaved struct {
	Group      int
	DateUpdate time.Time
	Schedule   string
}

type Lesson struct {
	//PrepodNameEnc string `json:"prepodNameEnc"`
	DayDate       string `json:"dayDate"`
	AudNum        string `json:"audNum"`
	DisciplName   string `json:"disciplName"`
	BuildNum      string `json:"buildNum"`
	OrgUnitName   string `json:"orgUnitName"`
	DayTime       string `json:"dayTime"`
	DayNum        string `json:"dayNum"`
	Potok         string `json:"potok"`
	PrepodName    string `json:"prepodName"`
	DisciplNum    string `json:"disciplNum"`
	OrgUnitId     string `json:"orgUnitId"`
	PrepodLogin   string `json:"prepodLogin"`
	DisciplType   string `json:"disciplType"`
	MarkedDeleted bool   `json:"markedDeleted"`
	//DisciplNameEnc string `json:"disciplNameEnc"`
}

type Schedule struct {
	Day3 []Lesson `json:"3,omitempty"`
	Day2 []Lesson `json:"2,omitempty"`
	Day1 []Lesson `json:"1,omitempty"`
	Day6 []Lesson `json:"6,omitempty"`
	Day5 []Lesson `json:"5,omitempty"`
	Day4 []Lesson `json:"4,omitempty"`
}

type Prepod struct {
	LessonType []string
	Name       string
	Lesson     string
}

type DeletedLessons struct {
	Id              int64  `json:"id,omitempty"`
	Groupid         int    `json:"groupid,omitempty"`
	Creator         int    `json:"creator,omitempty"`
	CreatorPlatform string `json:"creator_platform,omitempty"`
	LessonId        int    `json:"lesson_id,omitempty"`
	Date            pgtype.Date
	Uniqstring      string `json:"uniqstring"`
}

type DeletedLessonsMin struct {
	LessonId   int    `json:"lesson_id,omitempty"`
	Uniqstring string `json:"uniqstring"`
}

type LessonNew struct {
	DayDate string `json:"day_date"`
	DayNum  int    `json:"day_num"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}
