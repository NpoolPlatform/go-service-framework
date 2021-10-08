package envconf

import (
	"os"
	"strconv"
)

type envConf struct {
	environmentTarget string
	consulHost        string
	consulPort        int
}

var envConfInstance envConf

const (
	envValueUnknown = ""
)

func Init() {
	target := os.Getenv("ENV_ENVIRONMENT_TARGET")
	if target == envValueUnknown {
		panic("Variable ENV_ENVIRONMENT_TARGET is not set, it must be set in environment")
	}

	consulHost := os.Getenv("ENV_CONSUL_HOST")
	if consulHost == envValueUnknown {
		panic("Variable ENV_CONSUL_HOST is not set, it must be set in environment")
	}

	consulPortStr := os.Getenv("ENV_CONSUL_PORT")
	if consulPortStr == envValueUnknown {
		panic("Variable ENV_CONSUL_PORT is not set, it must be set in environment")
	}

	consulPort, err := strconv.Atoi(consulPortStr)
	if err != nil {
		panic("Variable ENV_CONSUL_PORT is invalid, it must be set as int in environment")
	}

	envConfInstance.environmentTarget = target
	envConfInstance.consulHost = consulHost
	envConfInstance.consulPort = consulPort
}

func EnvironmentTarget() string {
	return envConfInstance.environmentTarget
}

func ConsulHost() string {
	return envConfInstance.consulHost
}

func ConsulPort() int {
	return envConfInstance.consulPort
}
