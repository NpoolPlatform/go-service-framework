package envconf

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type envConf struct {
	EnvironmentTarget string
	ConsulHost        string
	ConsulPort        int
	ContainerID       string
	IPs               []string
}

const (
	envValueUnknown   = ""
	NotRunInContainer = "NOT-RUN-IN-CONTAINER"
)

var (
	inTesting = false
	EnvConf   = envConf{}
)

func Init() error {
	target := os.Getenv("ENV_ENVIRONMENT_TARGET")
	if target == envValueUnknown {
		return xerrors.Errorf("Variable ENV_ENVIRONMENT_TARGET is not set, it must be set in environment")
	}

	consulHost := os.Getenv("ENV_CONSUL_HOST")
	if consulHost == envValueUnknown {
		return xerrors.Errorf("Variable ENV_CONSUL_HOST is not set, it must be set in environment")
	}

	consulPortStr := os.Getenv("ENV_CONSUL_PORT")
	if consulPortStr == envValueUnknown {
		return xerrors.Errorf("Variable ENV_CONSUL_PORT is not set, it must be set in environment")
	}

	consulPort, err := strconv.Atoi(consulPortStr)
	if err != nil {
		return xerrors.Errorf("Variable ENV_CONSUL_PORT is invalid, it must be set as int in environment")
	}

	containerID, err := getContainerID()
	if err != nil {
		return xerrors.Errorf("Fail to get container ID: %v", err)
	}

	ips, err := getHostnames(true)
	if err != nil {
		return xerrors.Errorf("Fail to get host ip: %v", err)
	}

	EnvConf.EnvironmentTarget = target
	EnvConf.ConsulHost = consulHost
	EnvConf.ConsulPort = consulPort
	EnvConf.ContainerID = containerID
	EnvConf.IPs = ips

	return nil
}

func getContainerID() (string, error) {
	containerID := NotRunInContainer

	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		if os.IsNotExist(err) {
			return containerID, nil
		}
		return "", xerrors.Errorf("fail to read container id: %v", err)
	}
	defer file.Close()

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

func getHostnames(ip bool) ([]string, error) {
	var hostname []byte
	var err error

	if ip {
		hostname, err = exec.Command("hostname", "-I").Output()
		if err != nil {
			hostname, err = exec.Command("hostname", "-i").Output()
		}
	} else {
		hostname, err = exec.Command("hostname").Output()
	}

	// we ignore error of system which do not provide hostname
	if inTesting {
		return strings.Split(strings.TrimSpace(string(hostname)), " "), nil
	}

	return strings.Split(strings.TrimSpace(string(hostname)), " "), err
}
