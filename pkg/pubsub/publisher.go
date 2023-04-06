package pubsub

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/google/uuid"
)

func Publish(mid string, uid, rid, unid *uuid.UUID, body interface{}) error {
	amqpConfig, err := DurablePubSubConfig()
	if err != nil {
		return err
	}
	publisher, err := amqp.NewPublisher(
		*amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return err
	}

	byteMsg, err := json.Marshal(body)
	if err != nil {
		return err
	}

	sendMsg := Msg{
		MsgBase: MsgBase{
			MID:    mid,
			Sender: Sender(),
			RID:    rid,
			UnID:   unid,
		},
		Body: string(byteMsg),
	}

	sendByteMsg, err := json.Marshal(sendMsg)
	if err != nil {
		return err
	}

	_uid := uuid.NewString()
	if uid != nil {
		_uid = uid.String()
	}

	return publisher.Publish(
		GlobalPubsubTopic,
		message.NewMessage(
			_uid,
			sendByteMsg,
		),
	)
}
