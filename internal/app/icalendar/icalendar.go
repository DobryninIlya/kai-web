package icalendar

type ICalendar struct {
	Body []byte
}

func NewICalendar() *ICalendar {
	return &ICalendar{}
}
