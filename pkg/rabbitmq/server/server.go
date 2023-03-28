package server

import (
	"encoding/json"
	"github.com/google/uuid"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/common" //nolint

	"github.com/streadway/amqp"
)

type server struct {
	*rabbitmq.RabbitMQ
}

var myServer = &server{}

func Init() error {
	mq, err := rabbitmq.New(rabbitmq.MyServiceNameToVHost())
	if err != nil {
		return xerrors.Errorf("fail to create rabbitmq: %v", err)
	}

	myServer.RabbitMQ = mq

	return nil
}

func Deinit() {
	if myServer != nil {
		myServer.Destroy()
	}
}

func DeclareQueue(queueName string) error {
	return myServer.DeclareQueue(queueName)
}

func PublishToQueue(queueName string, msg interface{}) error {
	_, ok := myServer.Queues[queueName]
	if !ok {
		return xerrors.Errorf("queue '%v' is not declared, call DeclareQueue firstly", queueName)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return xerrors.Errorf("fail to marshal queue '%v' msg: %v", queueName, err)
	}

	return myServer.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "applition/json",
			DeliveryMode: amqp.Persistent,
			Body:         b,
		},
	)
}

func PublishToExchange(exChangeName string, msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return xerrors.Errorf("fail to marshal exchange '%v' msg: %v", exChangeName, err)
	}

	return myServer.Channel.Publish(
		exChangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "applition/json",
			DeliveryMode: amqp.Persistent,
			Body:         b,
			MessageId:    uuid.NewString(),
		},
	)
}
