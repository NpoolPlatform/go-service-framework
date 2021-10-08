package consul

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "127.0.0.1")
	os.Setenv("ENV_CONSUL_PORT", "8500")
}

func TestMain(m *testing.M) {
	command := exec.Command("consul", "agent", "-dev")
	go func() {
		if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); runByGithubAction || err != nil {
			return
		}

		_, err := command.Output()
		if err != nil {
			fmt.Printf("local consul environment is not prepared: %v\n", err)
		}
	}()

	time.Sleep(3 * time.Second)

	exitVal := m.Run()

	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); !runByGithubAction && err == nil {
		exec.Command("kill", "-9", fmt.Sprintf("%v", command.Process.Pid))
	}
	os.Exit(exitVal)
}

func TestNewConsulClient(t *testing.T) {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); runByGithubAction || err != nil {
		return
	}

	_, err := NewConsulClient()
	if err != nil {
		t.Errorf("fail to create consul client: %v", err)
	}
}

func TestRegisterService(t *testing.T) {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); runByGithubAction || err != nil {
		return
	}

	cli, err := NewConsulClient()
	if err != nil {
		t.Errorf("fail to create consul client: %v", err)
	}

	err = cli.RegisterService(RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1234,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	err = cli.RegisterService(RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1235,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	err = cli.RegisterService(RegisterInput{
		ID:   uuid.New(),
		Name: "unit-test-service",
		Tags: []string{"test", "unit-test"},
		Port: 1236,
	})
	if err != nil {
		t.Errorf("fail to register service: %v", err)
	}

	services, err := cli.QueryServices("unit-test-service")
	if err != nil {
		t.Errorf("fail to query services: %v", err)
	}

	if len(services) != 3*len(cli.envConf.IPs) {
		t.Errorf("service count is %v, expect %v", len(services), 3*len(cli.envConf.IPs))
	}
}
