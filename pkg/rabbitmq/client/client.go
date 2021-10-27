package client

import (
	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/common" //nolint

	"github.com/streadway/amqp"
)

type Client struct {
	mq *rabbitmq.RabbitMQ
}

func New(serviceName string) (*Client, error) {
	mq, err := rabbitmq.New(rabbitmq.ServiceNameToVHost(serviceName))
	if err != nil {
		return nil, xerrors.Errorf("fail to create rabbitmq: %v", err)
	}
	return &Client{
		mq: mq,
	}, nil
}

func (c *Client) Destroy() {
	if c.mq != nil {
		c.mq.Destroy()
	}
}

func (c *Client) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return c.mq.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
