package rabbitmq

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const" //nolint

	"github.com/streadway/amqp"
)

type client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

var myClient = client{}

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

	ch, err := conn.Channel()
	if err != nil {
		return xerrors.Errorf("fail to construct rabbitmq channel: %v", err)
	}

	myClient.conn = conn
	myClient.channel = ch

	return nil
}

func ServiceNameToVHost() string {
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)
	return strings.ReplaceAll(myServiceName, ".", "-")
}
