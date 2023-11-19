package model

type SiteUser struct {
	Login    string
	Password string
}

type SiteUserInfo struct {
	FirstName  string
	LastName   string
	MiddleName string
	Group      int `json:"group,omitempty"`
}
