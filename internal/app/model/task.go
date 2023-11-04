package model

import "github.com/jackc/pgx/pgtype"

type Task struct {
	ID           int         `json:"id"`
	Header       string      `json:"header"`
	Body         string      `json:"body"`
	CreateDate   pgtype.Date `json:"create_date"`
	DeadlineDate pgtype.Date `json:"deadline_date"`
	Author       string      `json:"author"`
	Attachments  []string    `json:"attachments"`
	Groupname    int         `json:"groupname,omitempty"`
}
