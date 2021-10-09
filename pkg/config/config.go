package config

import (
	"flag"
	"fmt"
	"math/rand"

	"golang.org/x/xerrors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/apollo.v0"

	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	KeyHostname    = "hostname"
	KeyHTTPPort    = "http_port"
	KeyGRPCPort    = "grpc_port"
	KeyHealthzPort = "healthz_port"
)

type Config struct {
	EnvConf *envconf.EnvConf
}

var inTesting = false

const (
	apolloServiceName = "apollo.npool.top"
)

func Init(configPath, appName string, consulCli *consul.Client) (*Config, error) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, xerrors.Errorf("fail to bind flags: %v", err)
	}

	viper.SetConfigName(fmt.Sprintf("%s.viper", appName))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(fmt.Sprintf("/etc/%v", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%v", appName))
	viper.AddConfigPath(".")

	// Following're must for every service
	// config:
	//   hostname: my-service.npool.top
	//   http_port: 32759
	//   grpc_port: 32789
	//   healthz_port: 32799
	//   appid: "89089012783789789719823798127398",
	//
	if err := viper.ReadInConfig(); err != nil {
		return nil, xerrors.Errorf("fail to init config: %v", err)
	}

	services, err := consulCli.QueryServices(apolloServiceName)
	if err != nil {
		return nil, xerrors.Errorf("fail to query apollo services: %v", err)
	}

	if len(services) == 0 {
		return nil, xerrors.Errorf("0 apollo services're found")
	}

	cfg := &Config{}

	cfg.EnvConf, err = envconf.NewEnvConf()
	if err != nil {
		return nil, xerrors.Errorf("fail to create environment configuration: %v", err)
	}

	targetIdx := rand.Intn(len(services))
	var service *consulapi.AgentService
	currentIdx := 0

	for _, srv := range services {
		if currentIdx == targetIdx {
			service = srv
			break
		}
	}

	if !inTesting {
		err = apollo.StartWithConf(&apollo.Conf{
			AppID:      viper.GetString("appid"),
			Cluster:    cfg.EnvConf.EnvironmentTarget,
			Namespaces: []string{viper.GetString("hostname")},
			IP:         fmt.Sprintf("http://%v:%v", service.Address, service.Port),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail to start apollo client: %v", err)
		}
	}

	return cfg, nil
}

func (cfg *Config) GetIntValue(key string) int {
	val, got := cfg.getLocalValue(key)
	if got {
		return val.(int)
	}
	return apollo.GetIntValue(key, -1)
}

func (cfg *Config) GetStringValue(key string) string {
	val, got := cfg.getLocalValue(key)
	if got {
		return val.(string)
	}
	return apollo.GetStringValueWithNameSpace(viper.GetString("hostname"), key, "")
}

func (cfg *Config) getLocalValue(key string) (interface{}, bool) {
	switch key {
	case KeyHostname:
		return viper.GetString(key), true
	case KeyHTTPPort:
		return viper.GetInt(key), true
	case KeyGRPCPort:
		return viper.GetInt(key), true
	case KeyHealthzPort:
		return viper.GetInt(key), true
	}
	return nil, false
}
