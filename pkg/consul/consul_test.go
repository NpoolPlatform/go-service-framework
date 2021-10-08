package consul

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "consul-server.kube-system.svc.cluster.local")
	os.Setenv("ENV_CONSUL_PORT", "8500")
}

func TestNewConsulClient(t *testing.T) {
	_, err := NewConsulClient()
	if err != nil {
		t.Errorf("fail to create consul client: %v", err)
	}
}
