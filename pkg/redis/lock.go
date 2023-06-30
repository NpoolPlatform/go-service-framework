package redis

import (
	"context"
	"time"

	"golang.org/x/xerrors"
)

func TryLock(key string, expire time.Duration) error {
	cli, err := GetClient()
	if err != nil {
		return xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp := cli.SetNX(ctx, key, 1, expire)
	locked, err := resp.Result()
	if err != nil {
		return xerrors.Errorf("fail lock %v: %v", key, err)
	}

	if !locked {
		return xerrors.Errorf("fail lock %v", key)
	}

	return nil
}

func Unlock(key string) error {
	cli, err := GetClient()
	if err != nil {
		return xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp := cli.Del(ctx, key)
	_, err = resp.Result()
	if err != nil {
		return xerrors.Errorf("fail unlock %v: %v", key, err)
	}

	return nil
}
