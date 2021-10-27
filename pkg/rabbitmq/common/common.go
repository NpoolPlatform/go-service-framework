package rabbitmq

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const" //nolint

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
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
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)

	if username == "" {
		return nil, xerrors.Errorf("invalid username")
	}
	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}
	if myServiceName == "" {
		return nil, xerrors.Errorf("invalid service name")
	}

	rsl := fmt.Sprintf("amqp://%v:%v@%v:%v/%v", username, password, service.Address, service.Port, vhost)
	conn, err := amqp.Dial(rsl)
	if err != nil {
		return nil, xerrors.Errorf("fail to create rabbitmq connection: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, xerrors.Errorf("fail to construct rabbitmq channel: %v", err)
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
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

func MyServiceNameToVHost() string {
	return ServiceNameToVHost(config.GetStringValueWithNameSpace("", config.KeyHostname))
}

func ServiceNameToVHost(serviceName string) string {
	return strings.ReplaceAll(serviceName, ".", "-")
}
