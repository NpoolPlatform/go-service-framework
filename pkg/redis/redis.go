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

const (
	keyUsername = "username"
	keyPassword = "password"
)

var myClient = Client{}

func Init() (*Client, error) {
	service, err := config.PeekService(constant.RedisServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query redis service: %v", err)
	}

	username := config.GetStringValueWithNameSpace(constant.RedisServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.RedisServiceName, keyPassword)

	if username == "" {
		return nil, xerrors.Errorf("invalid username")
	}
	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}

	myClient.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", service.Address, service.Port),
		Password: password,
		Username: username,
		DB:       0,
	})

	return &myClient, nil
}
