package publisher

import (
	"encoding/json"

	"github.com/NpoolPlatform/go-service-framework/pkg/eventbus"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

//nolint:deadcode
func Publisher(messageID string, msg interface{}) error {
	amqpConfig, err := eventbus.DurablePubSubConfig()
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

	byteMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	sendMsg := eventbus.Message{
		MessageID: messageID,
		UniqueID:  uuid.New(),
		Body:      byteMsg,
	}

	sendByteMsg, err := json.Marshal(sendMsg)
	if err != nil {
		return err
	}

	return publisher.Publish(eventbus.Topic, message.NewMessage(sendMsg.UniqueID.String(), sendByteMsg))
}
