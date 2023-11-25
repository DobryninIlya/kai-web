package model

import "github.com/jackc/pgx/pgtype"

type News struct {
	Id          int         `json:"id,omitempty"`
	Header      string      `json:"header,omitempty"`
	Description string      `json:"description,omitempty"`
	Body        string      `json:"body,omitempty"`
	Date        pgtype.Date `json:"date,omitempty"`
	Tag         string      `json:"tag,omitempty"`
	PreviewURL  string      `json:"preview_url,omitempty"`
	Author      int         `json:"author"`
	AuthorName  string      `json:"author"`
	AICorrect   bool        `json:"ai_correct"`
	Views       int         `json:"views"`
}
