package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const" //nolint
	"golang.org/x/xerrors"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queues  map[string]*amqp.Queue
}

const (
	keyUsername = "username"
	keyPassword = "password"
)

func New(vhost string) (*RabbitMQ, error) {
	service, err := config.PeekService(constant.RabbitMQServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query rabbitmq service: %v", err)
	}

	username := config.GetStringValueWithNameSpace(constant.RabbitMQServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.RabbitMQServiceName, keyPassword)

	if username == "" {
		return nil, xerrors.Errorf("invalid username")
	}
	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}

	rsl := fmt.Sprintf("amqp://%v:%v@%v:%v/%v", username, password, service.Address, service.Port, vhost)
	conn, err := amqp.Dial(rsl)
	if err != nil {
		return nil, xerrors.Errorf("fail to create rabbitmq %v connection: %v", rsl, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, xerrors.Errorf("fail to construct rabbitmq %v channel: %v", rsl, err)
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		Queues:  map[string]*amqp.Queue{},
	}, nil
}

func (mq *RabbitMQ) Destroy() {
	if mq.Channel != nil {
		mq.Channel.Close()
	}
	if mq.Conn != nil {
		mq.Conn.Close()
	}
}

func (mq *RabbitMQ) DeclareQueue(queueName string) error {
	queue, err := mq.Channel.QueueDeclare(
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

	mq.Queues[queueName] = &queue

	return nil
}

func MyServiceNameToVHost() string {
	return ServiceNameToVHost(config.GetStringValueWithNameSpace("", config.KeyHostname))
}

func ServiceNameToVHost(serviceName string) string {
	return strings.ReplaceAll(serviceName, ".", "-")
}
