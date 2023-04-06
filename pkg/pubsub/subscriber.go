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

type MsgHandler func(ctx context.Context, mid string, uid uuid.UUID, respToID *uuid.UUID, body string) error

func Subscribe(ctx context.Context, handler MsgHandler) error {
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

func processMsg(ctx context.Context, msg *message.Message, handler MsgHandler) {
	// We always need to ack, unless we're crashed or exit
	// Watermill set autoAck to false, and wait for ack synchronized after each msg,
	// so we need to ack at first glance we get the msg
	// https://www.rabbitmq.com/consumers.html#acknowledgement-timeout
	msg.Ack()

	msg1 := Msg{}
	err := json.Unmarshal(msg.Payload, &msg1)
	if err != nil {
		logger.Sugar().Errorw(
			"processMsg",
			"UUID", msg.UUID,
			"Metadata", msg.Metadata,
			"Payload", msg.Payload,
			"Error", err,
		)
		return
	}

	logger.Sugar().Infow(
		"processMsg",
		"MID", msg1.MID,
		"Sender", msg1.Sender,
		"UUID", msg.UUID,
		"Body", msg1.Body,
		"RID", msg1.RID,
	)

	err = handler(
		ctx,
		msg1.MID,
		uuid.MustParse(msg.UUID),
		msg1.RID,
		msg1.Body,
	)
	if err != nil {
		logger.Sugar().Errorw(
			"processMsg",
			"MID", msg1.MID,
			"Sender", msg1.Sender,
			"UUID", msg.UUID,
			"Body", msg1.Body,
			"RID", msg1.RID,
			"Error", err,
		)
	}
}

func process(ctx context.Context, messages <-chan *message.Message, handler MsgHandler) {
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				logger.Sugar().Warnw(
					"process",
					"State", "Closed",
				)
				return
			}
			processMsg(ctx, msg, handler)
		case <-ctx.Done():
			logger.Sugar().Warnw(
				"process",
				"State", "Done",
				"Error", ctx.Err(),
			)
			return
		}
	}
}
