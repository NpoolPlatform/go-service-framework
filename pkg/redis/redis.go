package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"github.com/go-redis/redis/v8"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/redis/const"
)

type Client struct {
	Client *redis.Client
}

const keyPassword = "password"

var (
	myClient = Client{}
	myMutex  sync.Mutex
	pingFail = false
)

func init() {
	ping()
}

func newClient() (*redis.Client, error) {
	service, err := config.PeekService(constant.RedisServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query redis service: %v", err)
	}

	password := config.GetStringValueWithNameSpace(constant.RedisServiceName, keyPassword)

	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", service.Address, service.Port),
		Password: password,
		DB:       0,
	})

	return client, nil
}

func GetClient() (*redis.Client, error) {
	myMutex.Lock()
	if !pingFail && myClient.Client != nil {
		cli := myClient.Client
		myMutex.Unlock()
		return cli, nil
	}

	if myClient.Client != nil {
		myClient.Client.Close()
	}
	cli, err := newClient()
	if err != nil {
		myMutex.Unlock()
		return nil, xerrors.Errorf("fail create redis client: %v", err)
	}

	pingFail = false
	myClient.Client = cli

	myMutex.Unlock()

	return cli, nil
}

func ping() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			myMutex.Lock()

			if myClient.Client == nil {
				cancel()
				myMutex.Unlock()
				continue
			}

			_, err := myClient.Client.Ping(ctx).Result()
			cancel()
			if err == nil {
				myMutex.Unlock()
				continue
			}
			pingFail = true
			myMutex.Unlock()
		}
	}()
}
