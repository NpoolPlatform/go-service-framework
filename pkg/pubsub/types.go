package pubsub

import (
	"github.com/google/uuid"
)

type MessageBase struct {
	MessageID   string
	UniqueID    uuid.UUID
	Sender      string
	RespondToID *uuid.UUID
}

type MessageResp struct {
	Code    int
	Message string
	Body    []byte
}

type Message struct {
	MessageBase
	Body []byte
}
