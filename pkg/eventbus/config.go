package eventbus

import (
	"fmt"
	"strings"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

const (
	keyUsername = "username"
	keyPassword = "password"
	Topic       = "global"
)

type Message struct {
	BusinessID string // TODO:Defined in message
	MessageID  uuid.UUID
	Body       []byte
}

func myServiceNameToVHost() string {
	return serviceNameToVHost(config.GetStringValueWithNameSpace("", config.KeyHostname))
}

func serviceNameToVHost(serviceName string) string {
	return strings.ReplaceAll(serviceName, ".", "-")
}

func DurablePubSubConfig() (*amqp.Config, error) {
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

	rsl := fmt.Sprintf("amqp://%v:%v@%v:%v/%v", username, password, service.Address, service.Port, myServiceNameToVHost())

	amqpConfig := amqp.NewDurablePubSubConfig(rsl, func(topic string) string {
		return "good-middleware1"
	})
	return &amqpConfig, nil
}
