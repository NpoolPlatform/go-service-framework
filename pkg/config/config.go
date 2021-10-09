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
	mysqlconst "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	KeyLogDir      = "log-dir"
	KeyAppID       = "appid"
	KeyHostname    = "hostname"
	KeyHTTPPort    = "http_port"
	KeyGRPCPort    = "grpc_port"
	KeyHealthzPort = "healthz_port"
)

var inTesting = false

const (
	apolloServiceName = "apollo.npool.top"
)

func Init(configPath, appName string) error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return xerrors.Errorf("fail to bind flags: %v", err)
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
		return xerrors.Errorf("fail to init config: %v", err)
	}

	service, err := PeekService(apolloServiceName)
	if err != nil {
		return xerrors.Errorf("fail to find a usable service: %v", err)
	}

	if !inTesting {
		err = apollo.StartWithConf(&apollo.Conf{
			AppID:      viper.GetString(KeyAppID),
			Cluster:    envconf.EnvConf.EnvironmentTarget,
			Namespaces: []string{viper.GetString(KeyHostname), mysqlconst.MysqlServiceName},
			IP:         fmt.Sprintf("http://%v:%v", service.Address, service.Port),
		})
		if err != nil {
			return xerrors.Errorf("fail to start apollo client: %v", err)
		}
	}

	return nil
}

func GetIntValue(key string) int {
	val, got := getLocalValue(key)
	if got {
		return val.(int)
	}
	return apollo.GetIntValue(key, -1)
}

func GetStringValueWithNameSpace(namespace, key string) string {
	val, got := getLocalValue(key)
	if got {
		return val.(string)
	}
	return apollo.GetStringValueWithNameSpace(namespace, key, "")
}

func getLocalValue(key string) (interface{}, bool) {
	switch key {
	case KeyLogDir:
		return viper.GetString(key), true
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

func PeekService(serviceName string) (*consulapi.AgentService, error) {
	services, err := consul.QueryServices(apolloServiceName)
	if err != nil {
		return nil, xerrors.Errorf("fail to query apollo services: %v", err)
	}

	targetIdx := rand.Intn(len(services))
	currentIdx := 0

	for _, srv := range services {
		if currentIdx == targetIdx {
			return srv, nil
		}
		currentIdx++
	}

	return nil, xerrors.Errorf("fail to find suitable service for %v, expect %v, total %v", serviceName, targetIdx, len(services))
}
