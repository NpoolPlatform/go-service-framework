package subscriber

import (
	"context"
	"encoding/json"

	"github.com/NpoolPlatform/go-service-framework/pkg/eventbus"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

//nolint:deadcode
func Subscriber(ctx context.Context, handler func(message eventbus.Message) error) error {
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
	go process(messages, handler)
	return nil
}

func process(messages <-chan *message.Message, handler func(message eventbus.Message) error) {
	for msg := range messages {
		msg1 := eventbus.Message{}
		err := json.Unmarshal(msg.Payload, &msg1)
		if err != nil {
			return
		}
		err = handler(msg1)
		if err != nil {
			return
		}
		msg.Ack()
	}
}
