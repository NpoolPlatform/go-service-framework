package publisher

import (
	"encoding/json"

	rabbitmq "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/common"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"golang.org/x/xerrors"
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

func PubMessage(exChangeName string, msg interface{}) error {
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
