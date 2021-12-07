package config

import (
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/xerrors"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-archaius/source/apollo"
	"github.com/go-chassis/openlog"

	"github.com/spf13/viper"

	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	KeyLogDir         = "logdir"
	KeyAppID          = "appid"
	KeyServiceID      = "serviceid"
	KeyHostname       = "hostname"
	KeyHTTPPort       = "http_port"
	KeyGRPCPort       = "grpc_port"
	KeyPrometheusPort = "prometheus_port"
	rootConfig        = "config"
)

var inTesting = false

const (
	apolloServiceName = "apollo-configservice.npool.top"
)

type Listener struct {
	Key string
}

func (l *Listener) Event(ev *event.Event) {
	openlog.Info(ev.Key)
	openlog.Info(fmt.Sprintf("%v\n", ev.Value))
	openlog.Info(ev.EventType)
}

// Init deps dependent other services
func Init(configPath, appName string, deps ...string) error {
	viper.SetConfigName(fmt.Sprintf("%s.viper", appName))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("./")
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

	appID := viper.GetStringMap(rootConfig)[KeyAppID].(string)         //nolint
	serviceID := viper.GetStringMap(rootConfig)[KeyServiceID].(string) //nolint
	myHostname := viper.GetStringMap(rootConfig)[KeyHostname].(string) //nolint
	logDir := viper.GetStringMap(rootConfig)[KeyLogDir].(string)       //nolint

	depServices := make([]string, len(deps)+1)
	for idx, dep := range deps {
		depServices[idx] = serviceNameToNamespace(dep)
	}
	depServices[len(depServices)-1] = serviceNameToNamespace(myHostname)
	namespaces := strings.Join(depServices, ",")

	fmt.Printf("cluster: %v\n", envconf.EnvConf.EnvironmentTarget)
	fmt.Printf("namespace: %v\n", namespaces)
	fmt.Printf("appid: %v\n", appID)
	fmt.Printf("serviceid: %v\n", serviceID)
	fmt.Printf("logdir: %v\n", logDir)

	if !inTesting {
		err := archaius.Init(
			archaius.WithRemoteSource(archaius.ApolloSource, &archaius.RemoteInfo{
				URL: fmt.Sprintf("http://%v:%v", service.Address, service.Port),
				DefaultDimension: map[string]string{
					apollo.AppID:         appID,
					apollo.NamespaceList: namespaces,
					apollo.Cluster:       envconf.EnvConf.EnvironmentTarget,
				},
			}))
		if err != nil {
			return xerrors.Errorf("fail to start apollo client: %v", err)
		}

		err = archaius.RegisterListener(&Listener{})
		if err != nil {
			return xerrors.Errorf("fail to register listener: %v", err)
		}
	}

	return nil
}

func GetIntValueWithNameSpace(namespace, key string) int {
	val, got := getLocalValue(key)
	if val != nil && got {
		return val.(int)
	}
	return archaius.GetInt(serviceNameKeyToApolloKey(serviceNameToNamespace(namespace), key), -1)
}

func GetStringValueWithNameSpace(namespace, key string) string {
	val, got := getLocalValue(key)
	if val != nil && got {
		return val.(string)
	}
	return archaius.GetString(serviceNameKeyToApolloKey(serviceNameToNamespace(namespace), key), "")
}

func getLocalValue(key string) (interface{}, bool) {
	switch key {
	case KeyLogDir:
		fallthrough //nolint
	case KeyHostname:
		fallthrough //nolint
	case KeyHTTPPort:
		fallthrough //nolint
	case KeyGRPCPort:
		fallthrough //nolint
	case KeyServiceID:
		fallthrough //nolint
	case KeyPrometheusPort:
		return viper.GetStringMap(rootConfig)[key], true
	}
	return nil, false
}

func PeekService(serviceName string, tags ...string) (*consulapi.AgentService, error) {
	services, err := consul.QueryServices(serviceName, tags...)
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

func ServiceNameToNamespace(serviceName string) string {
	return serviceNameToNamespace(serviceName)
}

func serviceNameToNamespace(serviceName string) string {
	return strings.ReplaceAll(serviceName, ".", "-")
}

func serviceNameKeyToApolloKey(serviceName, key string) string {
	return strings.Join([]string{serviceName, key}, ".")
}
