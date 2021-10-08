package envconf

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "consul-server.kube-system.svc.cluster.local")
	os.Setenv("ENV_CONSUL_PORT", "8500")
}

func TestNewEnvConf(t *testing.T) {
	conf, err := NewEnvConf()
	if err != nil {
		t.Errorf("fail to create env configuration: %v", err)
	}
	if conf.EnvironmentTarget != "development" {
		t.Errorf("target is %v, expect development", conf.EnvironmentTarget)
	}
	if conf.ConsulHost != "consul-server.kube-system.svc.cluster.local" {
		t.Errorf("consul host is %v, expect consul-server.kube-system.svc.cluster.local", conf.ConsulHost)
	}
	if conf.ConsulPort != 8500 {
		t.Errorf("consul port is %v, expect 8500", conf.ConsulPort)
	}
}

func TestGetContainerID(t *testing.T) {
	id, err := getContainerID()
	if err != nil {
		t.Errorf("fail to get container id: %v", err)
	}
	t.Logf("container id is %v", id)
}
