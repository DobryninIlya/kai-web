package model

import "time"

type ScheduleSaved struct {
	Group      int
	DateUpdate time.Time
	Schedule   string
}

type Lesson struct {
	PrepodNameEnc  string `json:"prepodNameEnc"`
	DayDate        string `json:"dayDate"`
	AudNum         string `json:"audNum"`
	DisciplName    string `json:"disciplName"`
	BuildNum       string `json:"buildNum"`
	OrgUnitName    string `json:"orgUnitName"`
	DayTime        string `json:"dayTime"`
	DayNum         string `json:"dayNum"`
	Potok          string `json:"potok"`
	PrepodName     string `json:"prepodName"`
	DisciplNum     string `json:"disciplNum"`
	OrgUnitId      string `json:"orgUnitId"`
	PrepodLogin    string `json:"prepodLogin"`
	DisciplType    string `json:"disciplType"`
	DisciplNameEnc string `json:"disciplNameEnc"`
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
