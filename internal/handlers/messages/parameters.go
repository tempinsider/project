package messages

import "time"

type ListRequest struct {
}

type Error struct {
	Error error `json:"ERROR"`
}

type ListResponse struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	SentAt            time.Time `json:"sent_at"`
	ExternalMessageID string    `json:"external_message_id"`
}
