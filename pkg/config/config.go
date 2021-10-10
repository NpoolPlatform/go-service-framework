package config

import (
	"fmt"
	"math/rand"
	"strconv"

	"golang.org/x/xerrors"

	"github.com/shima-park/agollo"

	"github.com/spf13/viper"

	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	mysqlconst "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	KeyLogDir      = "logdir"
	KeyAppID       = "appid"
	KeyHostname    = "hostname"
	KeyHTTPPort    = "http_port"
	KeyGRPCPort    = "grpc_port"
	KeyHealthzPort = "healthz_port"
)

type config struct {
	agollo.Agollo
}

var (
	inTesting = false
	myConfig  = config{}
)

const (
	apolloServiceName = "apollo.npool.top"
)

func Init(configPath, appName string) error {
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
	//   logdir: "/var/log"
	//
	if err := viper.ReadInConfig(); err != nil {
		return xerrors.Errorf("fail to init config: %v", err)
	}

	service, err := PeekService(apolloServiceName)
	if err != nil {
		return xerrors.Errorf("fail to find a usable service: %v", err)
	}

	if !inTesting {
		cli, err := agollo.New(
			fmt.Sprintf("http://%v:%v", service.Address, service.Port),
			viper.GetString(KeyAppID),
			agollo.PreloadNamespaces(viper.GetString(KeyHostname), mysqlconst.MysqlServiceName),
			agollo.AutoFetchOnCacheMiss(),
		)
		if err != nil {
			return xerrors.Errorf("fail to start apollo client: %v", err)
		}

		myConfig = config{
			Agollo: cli,
		}
	}

	return nil
}

func GetIntValueWithNameSpace(namespace, key string) int {
	val, got := getLocalValue(key)
	if got {
		return val.(int)
	}
	rval := myConfig.Get(key, agollo.WithNamespace(namespace))
	ival, err := strconv.ParseInt(rval, 10, 32)
	if err != nil {
		return -1
	}

	return int(ival)
}

func GetStringValueWithNameSpace(namespace, key string) string {
	val, got := getLocalValue(key)
	if got {
		return val.(string)
	}
	return myConfig.Get(key, agollo.WithNamespace(namespace))
}

func getLocalValue(key string) (interface{}, bool) {
	switch key {
	case KeyLogDir:
		return viper.GetStringMap("config")[key], true
	case KeyHostname:
		return viper.GetStringMap("config")[key], true
	case KeyHTTPPort:
		return viper.GetStringMap("config")[key], true
	case KeyGRPCPort:
		return viper.GetStringMap("config")[key], true
	case KeyHealthzPort:
		return viper.GetStringMap("config")[key], true
	}
	return nil, false
}

func PeekService(serviceName string) (*consulapi.AgentService, error) {
	services, err := consul.QueryServices(serviceName)
	if err != nil {
		return nil, xerrors.Errorf("fail to query apollo services: %v", err)
	}

	if len(services) == 0 {
		return nil, xerrors.Errorf("fail to find services of %v", serviceName)
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
