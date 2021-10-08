package config

import (
	"flag"
	"fmt"
	"math/rand"

	"golang.org/x/xerrors"

	"github.com/philchia/agollo/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

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
	agollo.Client
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

	cfg.Client = agollo.NewClient(&agollo.Conf{
		AppID:          viper.GetString("appid"),
		Cluster:        cfg.EnvConf.EnvironmentTarget,
		NameSpaceNames: []string{},
		MetaAddr:       fmt.Sprintf("%v:%v", service.Address, service.Port),
	})

	err = cfg.Start()
	if err != nil {
		return nil, xerrors.Errorf("fail to start apollo client: %v", err)
	}

	return cfg, nil
}

func (cfg *Config) Start() error {
	if !inTesting {
		return cfg.Client.Start()
	}
	return nil
}
