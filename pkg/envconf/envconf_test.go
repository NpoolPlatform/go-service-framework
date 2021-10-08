package envconf

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "consul-server.kube-system.svc.cluster.local")
	os.Setenv("ENV_CONSUL_PORT", "8500")
	Init()
}

func TestEnvironmentTarget(t *testing.T) {
	target := EnvironmentTarget()
	if target != "development" {
		t.Errorf("target is %v, expect development", target)
	}
}

func TestConsulHostAndPort(t *testing.T) {
	host := ConsulHost()
	if host != "consul-server.kube-system.svc.cluster.local" {
		t.Errorf("consul host is %v, expect consul-server.kube-system.svc.cluster.local", host)
	}

	port := ConsulPort()
	if port != 8500 {
		t.Errorf("consul port is %v, expect 8500", port)
	}
}
