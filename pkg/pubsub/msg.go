package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type MsgBase struct {
	MID    string
	UID    uuid.UUID
	Sender string
	RID    *uuid.UUID
	UnID   *uuid.UUID
}

type Resp struct {
	Code int
	Msg  string
}

type Msg struct {
	MsgBase
	Body string
	priv interface{}
}

func (msg *Msg) Nack() {
	msg.priv.(*message.Message).Nack()
}

func (msg *Msg) Ack() {
	msg.priv.(*message.Message).Ack()
}
