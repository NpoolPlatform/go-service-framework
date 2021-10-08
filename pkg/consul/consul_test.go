package consul

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "127.0.0.1")
	os.Setenv("ENV_CONSUL_PORT", "8500")
}

func TestMain(m *testing.M) {
	command := exec.Command("consul", "agent", "-dev")
	go func() {
		_, err := command.Output()
		if err != nil {
			fmt.Printf("local consul environment is not prepared: %v\n", err)
		}
	}()

	time.Sleep(3 * time.Second)

	exitVal := m.Run()

	exec.Command("kill", "-9", fmt.Sprintf("%v", command.Process.Pid))
	os.Exit(exitVal)
}

func TestNewConsulClient(t *testing.T) {
	_, err := NewConsulClient()
	if err != nil {
		t.Errorf("fail to create consul client: %v", err)
	}
}
