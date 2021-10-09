package mysql

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
)

type Client struct{}

func NewMysqlClient(cfg *config.Config, consulCli *consul.Client) (*Client, error) {
	return nil, nil
}
