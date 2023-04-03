package pubsub

import (
	"github.com/google/uuid"
)

type MessageBase struct {
	MessageID  string
	UniqueID   uuid.UUID
	Sender     string
	ResponseID string
}

type Message struct {
	MessageBase
	Body []byte
}
