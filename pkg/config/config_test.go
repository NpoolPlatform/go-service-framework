package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
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
}

func TestMain(m *testing.M) {
	err := os.MkdirAll(cfgDir, 0755) //nolint
	if err != nil {
		panic(fmt.Sprintf("fail to create dir %v: %v", cfgDir, err))
	}

	type Apollo struct {
		AppID          string   `yaml:"appid"`
		Cluster        string   `yaml:"cluster"`
		NameSpaceNames []string `yaml:"namespacenames"`
		MetaAddr       string   `yaml:"metaaddr"`
	}

	type T struct {
		Apollo Apollo `yaml:"apollo"`
	}

	apolloCfg, err := yaml.Marshal(T{
		Apollo: Apollo{
			AppID:          "123123123",
			Cluster:        "default",
			NameSpaceNames: []string{"123123123"},
			MetaAddr:       "localhost",
		},
	})
	if err != nil {
		panic(fmt.Sprintf("fail to marshal apollo config: %v", err))
	}

	err = ioutil.WriteFile(cfgFile, []byte(apolloCfg), 0755) //nolint
	if err != nil {
		panic(fmt.Sprintf("fail to write %v: %v", cfgFile, err))
	}

	inTesting = true

	exitVal := m.Run()

	os.Remove(cfgFile)
	os.Exit(exitVal)
}

func TestInit(t *testing.T) {
	_, err := Init(cfgDir, cfgName)
	if err != nil {
		t.Errorf("cannot init config: %v", err)
	}
}
