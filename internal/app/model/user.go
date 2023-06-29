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

// Validate ...
//func (u *User) Validate() error {
//	return validation.ValidateStruct(
//		u,
//		validation.Field(&u.Email, validation.Required, is.Email),
//		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
//	)
//}

//// BeforeCreate ...
//func (u *User) BeforeCreate() error {
//	if len(u.Password) > 0 {
//		enc, err := encryptString(u.Password)
//		if err != nil {
//			return err
//		}
//
//		u.EncryptedPassword = enc
//	}
//
//	return nil
//}
//
//// Sanitize ...
//func (u *User) Sanitize() {
//	u.Password = ""
//}
//
//// ComparePassword ...
//func (u *User) ComparePassword(password string) bool {
//	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
//}
//
//func encryptString(s string) (string, error) {
//	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
//	if err != nil {
//		return "", err
//	}
//
//	return string(b), nil
//}
