package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus int

const (
	MessageStatusOpen MessageStatus = iota
	MessageStatusSent
)

type Message struct {
	ID                uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	MessageContent    string
	PhoneNumber       string
	Status            MessageStatus
	ExternalMessageID *string
}
