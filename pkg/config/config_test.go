package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"

	"github.com/google/uuid"
)

const (
	cfgDir  = "./tmp"
	cfgName = "test"
)

var cfgFile = fmt.Sprintf("%s/%s.viper.yaml", cfgDir, cfgName)

func init() {
	os.Setenv("ENV_ENVIRONMENT_TARGET", "development")
	os.Setenv("ENV_CONSUL_HOST", "consul-server.kube-system.svc.cluster.local")
	os.Setenv("ENV_CONSUL_PORT", "8500")

	err := envconf.Init()
	if err != nil {
		fmt.Printf("fail to init environment config: %v", err)
	}

	err = consul.Init()
	if err != nil {
		panic(fmt.Sprintf("fail to init consul: %v", err))
	}
}

func TestMain(m *testing.M) {
	err := os.MkdirAll(cfgDir, 0755) //nolint
	if err != nil {
		panic(fmt.Sprintf("fail to create dir %v: %v", cfgDir, err))
	}

	type Config struct {
		AppID       string `yaml:"appid"`
		Hostname    string `yaml:"hostname"`
		HTTPPort    int    `yaml:"http_port"`
		GRPCPort    int    `yaml:"grpc_port"`
		HealthzPort int    `yaml:"healthz_port"`
		LogDir      string `yaml:"logdir"`
	}

	type T struct {
		Config Config `yaml:"config"`
	}

	config, err := yaml.Marshal(T{
		Config: Config{
			AppID:       "123123123123123123",
			Hostname:    "config-unit-test-service.npool.top",
			HTTPPort:    32578,
			GRPCPort:    32580,
			HealthzPort: 32582,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("fail to marshal service config: %v", err))
	}

	err = ioutil.WriteFile(cfgFile, []byte(config), 0755) //nolint
	if err != nil {
		panic(fmt.Sprintf("fail to write %v: %v", cfgFile, err))
	}

	inTesting = true

	exitVal := m.Run()

	os.Remove(cfgFile)
	os.Exit(exitVal)
}

func TestInit(t *testing.T) {
	if os.Getenv("RUN_BY_GITHUB_ACTION") == "true" {
		return
	}

	id := uuid.New()
	err := consul.RegisterService(true, consul.RegisterInput{
		ID:   id.String(),
		Name: "apollo.npool.top",
		Tags: []string{"apollo", "unit-test"},
		Port: 1235,
	})
	if err != nil {
		t.Errorf("fail to register apollo service: %v", err)
	}

	err = Init(cfgDir, cfgName)
	if err != nil {
		t.Errorf("cannot init config: %v", err)
	}

	err = consul.DeregisterService(id)
	if err != nil {
		t.Errorf("fail to deregister apollo service: %v", err)
	}
}
