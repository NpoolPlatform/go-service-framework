package envconf

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type EnvConf struct {
	EnvironmentTarget string
	ConsulHost        string
	ConsulPort        int
	ContainerID       string
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

func getContainerID() (string, error) {
	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return "", xerrors.Errorf("fail to read container id: %v", err)
	}
	defer file.Close()

	containerID := "NOT-RUN-IN-CONTAINER"

	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return "", xerrors.Errorf("fail to read cgroup file: %v", err)
		}

		strs := strings.Split(line, ":")
		if len(strs) < 3 {
			continue
		}

		if !strings.HasPrefix(strs[2], "/docker/") {
			continue
		}

		containerID = strings.ReplaceAll(strs[2], "/docker/", "")
		break
	}

	return containerID, nil
}
