package model

type ApiClient struct {
	Id         int    `json:"id,omitempty"`
	DeviceId   string `json:"device_id,omitempty"`
	DeviceTag  string `json:"device_tag,omitempty"`
	Token      string `json:"token,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
