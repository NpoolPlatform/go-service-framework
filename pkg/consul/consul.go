package consul

import (
	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"
)

type Client struct{}

func NewConsulClient() (*Client, error) {
	_, err := envconf.NewEnvConf()
	if err != nil {
		return nil, xerrors.Errorf("fail to create environment configuration: %v", err)
	}

	return nil, nil
}
