package model

import (
	"github.com/jackc/pgx/pgtype"
)

type ApiClient struct {
	UID        string      `json:"uid,omitempty"`
	DeviceTag  string      `json:"device_tag,omitempty"`
	Token      string      `json:"token,omitempty"`
	CreateDate pgtype.Date `json:"create_date,omitempty"`
}
