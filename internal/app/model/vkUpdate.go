package model

type Confirmation struct {
	Type    string `json:"type,omitempty"`
	GroupID int    `json:"group_id,omitempty"`
}

type VKUpdate struct {
	Confirmation
}
