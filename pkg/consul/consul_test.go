package consul

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "consul-server.kube-system.svc.cluster.local")
	os.Setenv("ENV_CONSUL_PORT", "8500")

	err := envconf.Init()
	if err != nil {
		fmt.Printf("fail to init environment config: %v", err)
	}
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestNewConsulClient(t *testing.T) {
	if runByGithubAction, _ := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); runByGithubAction { //nolint
		return
	}

	err := Init()
	if err != nil {
		t.Errorf("fail to create consul client: %v", err)
	}
}

func TestRegisterService(t *testing.T) {
	if runByGithubAction, _ := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); runByGithubAction { //nolint
		return
	}

	err := Init()
	if err != nil {
		t.Errorf("fail to create consul client: %v", err)
	}

	err = RegisterService(true, RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1234,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	err = RegisterService(true, RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1235,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	err = RegisterService(true, RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1236,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	services, err := QueryServices("unit-test-service")
	if err != nil {
		t.Errorf("fail to query services: %v", err)
	}

	if len(services) != 3*len(envconf.EnvConf.IPs) {
		t.Errorf("service count is %v, expect %v", len(services), 3*len(envconf.EnvConf.IPs))
	}
}
