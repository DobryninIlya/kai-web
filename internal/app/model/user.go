package model

import (
	"time"
)

// User ...
type User struct {
	ID           int       `json:"id_vk"`
	Name         string    `json:"name"`
	Group        int       `json:"group"`
	Distribution int8      `json:"distribution"`
	Admlevel     int8      `json:"admlevel"`
	GroupReal    int       `json:"groupreal"`
	DateChanged  time.Time `json:"dateChanged"`
	Balance      int       `json:"balance"`
	Distr        int       `json:"distr"`
	Warn         int8      `json:"warn"`
	Expiration   time.Time `json:"expiration"`
	BanHistory   int8      `json:"banHistory"`
	IsChecked    int8      `json:"isChecked"`
	Role         int8      `json:"role"`
	Login        string    `json:"login"`
	PotokLecture bool      `json:"potokLecture"`
	HasOwnShed   bool      `json:"hasOwnShed"`
	Affiliate    bool      `json:"affiliate"`
}

type RegistrationData struct {
	Role          int    `json:"id"`
	Identificator string `json:"Identificator"`
	VkId          int    `json:"vk_id"`
}
