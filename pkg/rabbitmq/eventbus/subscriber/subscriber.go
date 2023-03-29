package publisher

import (
	rabbitmq "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/common"
	"github.com/streadway/amqp"

	msgcli "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/client"
)

const exchange = "global"

var myClient = &msgcli.Client{}

func Init() error {
	var err error
	myClient, err = msgcli.New(rabbitmq.MyServiceNameToVHost())
	if err != nil {
		return err
	}
	err = myClient.DeclareSub(rabbitmq.MyServiceNameToVHost(), exchange)
	if err != nil {
		return err
	}
	return nil
}

func SubMessage() (<-chan amqp.Delivery, error) {
	return myClient.Consume(rabbitmq.MyServiceNameToVHost())
}
