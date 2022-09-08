package redis

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/xerrors"

	"github.com/go-redis/redis/v8"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/redis/const"
)

const keyPassword = "password"

var (
	redisClient *redis.Client
	poolSize    = 50
	lk          sync.RWMutex

	ErrRedisClientNotInit = errors.New("redis client not init")
)

func newClient() (*redis.Client, error) {
	lk.Lock()
	defer lk.Unlock()

	// double read
	if redisClient != nil {
		return redisClient, nil
	}

	service, err := config.PeekService(constant.RedisServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query redis service: %v", err)
	}

	password := config.GetStringValueWithNameSpace(constant.RedisServiceName, keyPassword)
	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", service.Address, service.Port),
		Password: password,
		DB:       0,
		PoolSize: poolSize,
	})

	return redisClient, nil
}

func GetClient() (*redis.Client, error) {
	lk.RLock()
	if redisClient != nil {
		_redisClient := redisClient
		lk.RUnlock()
		return _redisClient, nil
	}
	lk.RUnlock()

	var err error
	redisClient, err = newClient()
	return redisClient, err
}

func Close() error {
	lk.Lock()
	defer lk.Unlock()

	if redisClient != nil {
		redisClient.Close()
		redisClient = nil
	}

	return nil
}
