package redis

import (
	"context"
	"time"

	"golang.org/x/xerrors"
)

func Set(key string, value interface{}, expire time.Duration) error {
	cli, err := GetClient()
	if err != nil {
		return xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = cli.Set(ctx, key, value, expire).Err()
	if err != nil {
		return xerrors.Errorf("fail set key %v: %v", key, err)
	}

	return nil
}

func Get(key string) (interface{}, error) {
	cli, err := GetClient()
	if err != nil {
		return nil, xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	v, err := cli.Get(ctx, key).Result()
	if err != nil {
		return nil, xerrors.Errorf("fail get key %v: %v", key, err)
	}

	return v, nil
}

func Del(key string) error {
	cli, err := GetClient()
	if err != nil {
		return xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = cli.Del(ctx, key).Err()
	if err != nil {
		return xerrors.Errorf("fail del key %v: %v", key, err)
	}

	return nil
}
