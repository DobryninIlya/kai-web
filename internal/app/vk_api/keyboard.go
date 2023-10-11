package vk_api

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Inline  bool       `json:"inline"`
	Buttons [][]Button `json:"buttons"`
}

type Button interface {
}

type Action struct {
	Type string `json:"type"`
}

type TextButton struct {
	Action ActionText `json:"action"`
}

type ActionText struct {
	Type    string `json:"type"`
	Label   string `json:"label"`
	Payload string `json:"payload"`
}

type LinkButton struct {
	Action ActionLink `json:"action"`
}

type ActionLink struct {
	Type    string `json:"type"`
	Label   string `json:"label"`
	Link    string `json:"link"`
	Payload string `json:"payload"`
}

func GetInlineLinkButtonVK(Type, Label, Link, Payload string) Keyboard {
	return Keyboard{
		OneTime: false,
		Inline:  true,
		Buttons: [][]Button{
			[]Button{
				LinkButton{
					Action: ActionLink{
						Type:    Type,
						Label:   Label,
						Link:    Link,
						Payload: Payload,
					},
				},
			},
		},
	}
}
