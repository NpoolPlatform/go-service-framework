package subscriber

import (
	"context"
	"encoding/json"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/eventbus"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

//nolint:deadcode
func Subscriber(
	ctx context.Context,
	handler func(ctx context.Context, messageID string, UniqueID uuid.UUID, body []byte) error,
) error {
	amqpConfig, err := eventbus.DurablePubSubConfig()
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
	messages, err := subscriber.Subscribe(ctx, eventbus.Topic)
	if err != nil {
		return err
	}
	go process(ctx, messages, handler)
	return nil
}

func process(
	ctx context.Context,
	messages <-chan *message.Message,
	handler func(ctx context.Context, messageID string, UniqueID uuid.UUID, body []byte) error,
) {
	for msg := range messages {
		msg1 := eventbus.Message{}
		err := json.Unmarshal(msg.Payload, &msg1)
		if err != nil {
			return
		}
		for i := 0; i < 3; i++ {
			err = handler(ctx, msg1.MessageID, msg1.UniqueID, msg1.Body)
			if err == nil {
				msg.Ack()
				return
			}
			time.Sleep(time.Second * 5)
		}
		// TODO:send alarm messages
		logger.Sugar().Errorf("fail handler message id:%v,error:%v", msg1.MessageID, err)
	}
}
