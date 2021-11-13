package redis

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/go-redis/redis/v8"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/redis/const" //nolint
)

type Client struct {
	Client *redis.Client
}

const keyPassword = "password"

var myClient = Client{}

func Init() (*Client, error) {
	service, err := config.PeekService(constant.RedisServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query redis service: %v", err)
	}

	password := config.GetStringValueWithNameSpace(constant.RedisServiceName, keyPassword)

	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}

	myClient.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", service.Address, service.Port),
		Password: password,
		DB:       0,
	})

	return &myClient, nil
}
