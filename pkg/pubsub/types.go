package pubsub

import (
	"github.com/google/uuid"
)

type MsgBase struct {
	MID    string
	Sender string
	RID    *uuid.UUID
}

type Resp struct {
	Code int
	Msg  string
	Body string
}

type Msg struct {
	MsgBase
	Body string
}
