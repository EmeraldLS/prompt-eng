package types

import "time"

type ChatMessage struct {
	Prompt   string    `json:"prompt" validate:"required"`
	ChatTime time.Time `json:"chat_time"`
}

type ChatMessageResponse struct {
	Answer string `json:"answer"`
}
