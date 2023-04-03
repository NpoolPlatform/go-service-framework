package pubsub

import (
	"context"
	"encoding/json"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/google/uuid"
)

type MsgHandler func(ctx context.Context, messageID, sender string, uniqueID uuid.UUID, body []byte, respondToID *uuid.UUID) error

func Subscrib(ctx context.Context, handler MsgHandler) error {
	amqpConfig, err := DurablePubSubConfig()
	if err != nil {
		return err
	}

	subscriber, err := amqp.NewSubscriber(
		*amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return err
	}

	messages, err := subscriber.Subscribe(ctx, GlobalPubsubTopic)
	if err != nil {
		return err
	}

	go process(ctx, messages, handler)

	return nil
}

func process(ctx context.Context, messages <-chan *message.Message, handler MsgHandler) {
	for msg := range messages {
		msg1 := Message{}
		err := json.Unmarshal(msg.Payload, &msg1)
		if err != nil {
			logger.Sugar().Errorw("process", "Error", err)
			continue
		}

		logger.Sugar().Infow(
			"process",
			"MessageID", msg1.MessageID,
			"Sender", msg1.Sender,
			"UniqueID", msg1.UniqueID,
			"Body", string(msg1.Body),
			"ResponseToID", msg1.RespondToID,
		)

		err = handler(ctx, msg1.MessageID, msg1.Sender, msg1.UniqueID, msg1.Body, msg1.RespondToID)
		if err != nil {
			logger.Sugar().Errorw(
				"process",
				"MessageID", msg1.MessageID,
				"Sender", msg1.Sender,
				"UniqueID", msg1.UniqueID,
				"ResponseToID", msg1.RespondToID,
				"Error", err,
			)
			continue
		}
		msg.Ack()
	}
}
