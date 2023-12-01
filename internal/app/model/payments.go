package model

import "github.com/jackc/pgx/pgtype"

type Transaction struct {
	UID      string `json:"uid"`
	ExtID    string `json:"ext_id"`
	ClientID string `json:"client_id"`
	Date     pgtype.Date
	Type     string `json:"type"`
	Ended    bool   `json:"ended"`
}
