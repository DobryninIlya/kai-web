package model

type UploadServerAnswer struct {
	Response struct {
		AlbumID   int    `json:"album_id"`
		UploadURL string `json:"upload_url"`
		UserID    int    `json:"user_id"`
		GroupID   int    `json:"group_id"`
	} `json:"response"`
}
