package model

import "time"

type Verification struct {
	ID         int       `json:"id"`
	DateUpdate time.Time `json:"date_update"`
	Faculty    int       `json:"faculty"`
	Course     int       `json:"course"`
	Group      int       `json:"group"`
	IDcard     int       `json:"id_card"`
	Groupname  int       `json:"groupname"`
	Student    string    `json:"student"`
}

type VerificationParams struct {
	Faculty   int `json:"faculty"`
	Course    int `json:"course"`
	Group     int `json:"group"`
	Student   int `json:"student"`
	ID        int `json:"id"`
	Groupname int `json:"groupname"`
}

type UserVerification struct {
	Id         int
	DateUpdate time.Time
	Faculty    int
	Course     int
	GroupId    int
	Idcard     int
	Groupname  int
	Studentid  int
}
