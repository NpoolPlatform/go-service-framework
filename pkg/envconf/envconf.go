package envconf

import (
	"os"
	"strconv"

	"golang.org/x/xerrors"
)

type EnvConf struct {
	EnvironmentTarget string
	ConsulHost        string
	ConsulPort        int
}

const (
	envValueUnknown = ""
)

func NewEnvConf() (*EnvConf, error) {
	target := os.Getenv("ENV_ENVIRONMENT_TARGET")
	if target == envValueUnknown {
		return nil, xerrors.Errorf("Variable ENV_ENVIRONMENT_TARGET is not set, it must be set in environment")
	}

	consulHost := os.Getenv("ENV_CONSUL_HOST")
	if consulHost == envValueUnknown {
		return nil, xerrors.Errorf("Variable ENV_CONSUL_HOST is not set, it must be set in environment")
	}

	consulPortStr := os.Getenv("ENV_CONSUL_PORT")
	if consulPortStr == envValueUnknown {
		return nil, xerrors.Errorf("Variable ENV_CONSUL_PORT is not set, it must be set in environment")
	}

	consulPort, err := strconv.Atoi(consulPortStr)
	if err != nil {
		return nil, xerrors.Errorf("Variable ENV_CONSUL_PORT is invalid, it must be set as int in environment")
	}

	return &EnvConf{
		EnvironmentTarget: target,
		ConsulHost:        consulHost,
		ConsulPort:        consulPort,
	}, nil
}
