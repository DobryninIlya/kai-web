package model

type Confirmation struct {
	GroupID int `json:"group_id,omitempty"`
}

//type Object interface {
//}

type PhotoSizes struct {
	Height int    `json:"height"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type Attachment struct {
	Type  string `json:"type"`
	Photo struct {
		AlbumId      int          `json:"album_id"`
		Date         int          `json:"date"`
		Id           int          `json:"id"`
		OwnerId      int          `json:"owner_id"`
		AccessKey    string       `json:"access_key"`
		PostId       int          `json:"post_id"`
		Sizes        []PhotoSizes `json:"sizes"`
		Text         string       `json:"text"`
		UserId       int          `json:"user_id"`
		WebViewToken string       `json:"web_view_token"`
		HasTags      bool         `json:"has_tags"`
	}
}

type WallPost struct {
	Object struct {
		InnerType string `json:"inner_type"`
		CanEdit   int    `json:"can_edit"`
		CreatedBy int    `json:"created_by"`
		CanDelete int    `json:"can_delete"`
		Donut     struct {
			IsDonut bool `json:"is_donut"`
		} `json:"donut"`
		Comments struct {
			Count int `json:"count"`
		} `json:"comments"`
		MarkedAsAds                 int          `json:"marked_as_ads"`
		ShortTextRate               float64      `json:"short_text_rate"`
		CompactAttachmentsBeforeCut int          `json:"compact_attachments_before_cut"`
		Hash                        string       `json:"hash"`
		Type                        string       `json:"type"`
		Attachments                 []Attachment `json:"attachments"`
		AttachmentsMeta             struct {
			PrimaryMode string `json:"primary_mode"`
		} `json:"attachments_meta"`
		Date       int    `json:"date"`
		FromId     int    `json:"from_id"`
		Id         int    `json:"id"`
		IsFavorite bool   `json:"is_favorite"`
		OwnerId    int    `json:"owner_id"`
		PostType   string `json:"post_type"`
		Text       string `json:"text"`
	} `json:"object"`
}

type WallPostObject struct {
	WallPost
}

type VKUpdate struct {
	Type    string `json:"type,omitempty"`
	EventID string `json:"event_id,omitempty"`
	Version string `json:"v,omitempty"`
	GroupId int    `json:"group_id,omitempty"`
	WallPostObject
	Confirmation
}
