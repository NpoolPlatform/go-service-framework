package main

import (
	"encoding/json"

	"github.com/NpoolPlatform/go-service-framework/pkg/eventbus"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

//nolint:deadcode
func Publisher(businessID string, msg []byte) error {
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

	msg1 := eventbus.Message{
		MessageID:  uuid.New(),
		BusinessID: businessID,
		Body:       msg,
	}

	byteMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return publisher.Publish(eventbus.Topic, message.NewMessage(msg1.MessageID.String(), byteMsg))
}
