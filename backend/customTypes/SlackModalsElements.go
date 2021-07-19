package customTypes

import "github.com/slack-go/slack"

type CustomElements struct {
	Emoji bool   `json:"emoji"`
	Text  string `json:"text"`
	Type  string `json:"type"`
}

type CustomContextBlock struct {
	Type     string           `json:"type"`
	BlockID  string           `json:"block_id"`
	Elements []CustomElements `json:"elements"`
}

type SlackEventCallback struct {
	ResponseAction string                 `json:"response_action"`
	View           slack.ModalViewRequest `json:"view"`
}
