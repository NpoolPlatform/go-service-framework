package consul

import (
	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	*api.Client
}

func NewConsulClient() (*Client, error) {
	envConf, err := envconf.NewEnvConf()
	if err != nil {
		return nil, xerrors.Errorf("fail to create environment configuration: %v", err)
	}

	config := api.DefaultConfig()
	config.Address = envConf.ConsulHost
	client, err := api.NewClient(config)
	if err != nil {
		return nil, xerrors.Errorf("fail to create consul client: %v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
