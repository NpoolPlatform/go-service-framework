package consul

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

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
	command := exec.Command("consul", "agent", "-dev")
	go func() {
		if runByGithubAction, _ := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); runByGithubAction { //nolint
			return
		}

		_, err := command.Output()
		if err != nil {
			fmt.Printf("local consul environment is not prepared: %v\n", err)
		}
	}()

	time.Sleep(3 * time.Second)

	exitVal := m.Run()

	if runByGithubAction, _ := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); !runByGithubAction { //nolint
		exec.Command("kill", "-9", fmt.Sprintf("%v", command.Process.Pid))
	}
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

	err = RegisterService(RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1234,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	err = RegisterService(RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1235,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	err = RegisterService(RegisterInput{
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
