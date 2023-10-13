package model

import "github.com/jackc/pgx/pgtype"

type News struct {
	Id          int         `json:"id"`
	Header      string      `json:"header"`
	Description string      `json:"description"`
	Body        string      `json:"body"`
	Date        pgtype.Date `json:"date"`
}
