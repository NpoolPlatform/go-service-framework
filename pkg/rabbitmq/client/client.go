package client

import (
	rabbitmq "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/common" //nolint
	"github.com/streadway/amqp"
	"golang.org/x/xerrors"
)

type Client struct {
	*rabbitmq.RabbitMQ
}

func New(serviceName string) (*Client, error) {
	mq, err := rabbitmq.New(rabbitmq.ServiceNameToVHost(serviceName))
	if err != nil {
		return nil, xerrors.Errorf("fail to create rabbitmq: %v", err)
	}
	return &Client{
		RabbitMQ: mq,
	}, nil
}

func (c *Client) Destroy() {
	c.RabbitMQ.Destroy()
}

func (c *Client) Consume(queueName string) (<-chan amqp.Delivery, error) {
	_, ok := c.Queues[queueName]
	if !ok {
		return nil, xerrors.Errorf("queue '%v' is not declared, call DeclareQueue firstly", queueName)
	}

	return c.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
