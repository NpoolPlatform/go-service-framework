package pubsub

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func Publish(messageID string, respondToID *uuid.UUID, body interface{}) error {
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

	sendMsg := Message{
		MessageBase: MessageBase{
			MessageID:   messageID,
			UniqueID:    uuid.New(),
			Sender:      Sender(),
			RespondToID: respondToID,
		},
		Body: byteMsg,
	}

	sendByteMsg, err := json.Marshal(sendMsg)
	if err != nil {
		return err
	}

	return publisher.Publish(
		GlobalPubsubTopic,
		message.NewMessage(
			sendMsg.UniqueID.String(),
			sendByteMsg,
		),
	)
}
