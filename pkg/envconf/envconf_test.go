package envconf

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "consul-server.kube-system.svc.cluster.local")
	os.Setenv("ENV_CONSUL_PORT", "8500")
	inTesting = true
}

func TestNewEnvConf(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("fail to create env configuration: %v", err)
	}
	if EnvConf.EnvironmentTarget != "development" {
		t.Errorf("target is %v, expect development", EnvConf.EnvironmentTarget)
	}
	if EnvConf.ConsulHost != "consul-server.kube-system.svc.cluster.local" {
		t.Errorf("consul host is %v, expect consul-server.kube-system.svc.cluster.local", EnvConf.ConsulHost)
	}
	if EnvConf.ConsulPort != 8500 {
		t.Errorf("consul port is %v, expect 8500", EnvConf.ConsulPort)
	}
}

func TestGetContainerID(t *testing.T) {
	id, err := getContainerID()
	if err != nil {
		t.Errorf("fail to get container id: %v", err)
	}
	t.Logf("container id is %v", id)
}

func TestGetHostname(t *testing.T) {
	if os.Getenv("RUN_BY_GITHUB_ACTION") == "true" {
		return
	}

	hosts, err := getHostnames(true)
	if err != nil {
		t.Errorf("fail to get hostname with ip: %v", err)
	}
	t.Logf("success to get %v hostname without ip: %v", len(hosts), hosts)

	_, err = getHostnames(false)
	if err != nil {
		t.Errorf("fail to get hostname without ip: %v", err)
	}
}
