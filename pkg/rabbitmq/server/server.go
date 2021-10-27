package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const" //nolint

	"github.com/streadway/amqp"
)

type server struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]*amqp.Queue
}

var myServer = server{
	queues: map[string]*amqp.Queue{},
}

const (
	keyUsername = "username"
	keyPassword = "password"
)

func Init() error {
	service, err := config.PeekService(constant.RabbitMQServiceName)
	if err != nil {
		return xerrors.Errorf("Fail to query rabbitmq service: %v", err)
	}

	username := config.GetStringValueWithNameSpace(constant.RabbitMQServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.RabbitMQServiceName, keyPassword)
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)

	if username == "" {
		return xerrors.Errorf("invalid username")
	}
	if password == "" {
		return xerrors.Errorf("invalid password")
	}
	if myServiceName == "" {
		return xerrors.Errorf("invalid service name")
	}

	vhost := ServiceNameToVHost()

	rsl := fmt.Sprintf("amqp://%v:%v@%v:%v/%v", username, password, service.Address, service.Port, vhost)
	conn, err := amqp.Dial(rsl)
	if err != nil {
		return xerrors.Errorf("fail to create rabbitmq connection: %v", err)
	}

	myServer.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return xerrors.Errorf("fail to construct rabbitmq channel: %v", err)
	}

	myServer.channel = ch

	return nil
}

func Deinit() {
	if myServer.channel != nil {
		myServer.channel.Close()
	}
	if myServer.conn != nil {
		myServer.conn.Close()
	}
}

func ServiceNameToVHost() string {
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)
	return strings.ReplaceAll(myServiceName, ".", "-")
}

func DeclareQueue(queueName string) error {
	queue, err := myServer.channel.QueueDeclare(
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

	return myServer.channel.Publish(
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
