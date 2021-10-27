package server

import (
	"encoding/json"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/common" //nolint

	"github.com/streadway/amqp"
)

type server struct {
	mq     *rabbitmq.RabbitMQ
	queues map[string]*amqp.Queue
}

var myServer = server{
	queues: map[string]*amqp.Queue{},
}

func Init() error {
	mq, err := rabbitmq.New(rabbitmq.MyServiceNameToVHost())
	if err != nil {
		return xerrors.Errorf("fail to create rabbitmq: %v", err)
	}

	myServer.mq = mq

	return nil
}

func Deinit() {
	if myServer.mq != nil {
		myServer.mq.Destroy()
	}
}

func DeclareQueue(queueName string) error {
	queue, err := myServer.mq.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return xerrors.Errorf("fail to construct rabbitmq queue %v: %v", queueName, err)
	}

	myServer.queues[queueName] = &queue

	return nil
}

func PublishToQueue(queueName string, msg interface{}) error {
	_, ok := myServer.queues[queueName]
	if !ok {
		return xerrors.Errorf("queue '%v' is not declared, call DeclareQueue firstly", queueName)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return xerrors.Errorf("fail to marshal queue '%v' msg: %v", queueName, err)
	}

	return myServer.mq.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
}
