package consul

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"

	"github.com/hashicorp/consul/api"
)

type client struct {
	*api.Client
}

var myClient *client

func Init() error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%v:%v", envconf.EnvConf.ConsulHost, envconf.EnvConf.ConsulPort)
	cli, err := api.NewClient(config)
	if err != nil {
		return xerrors.Errorf("fail to create consul client: %v", err)
	}

	myClient = &client{
		Client: cli,
	}

	return nil
}

// IP is parsed from package envconf
type RegisterInput struct {
	ID          uuid.UUID
	Name        string
	Tags        []string
	Port        int
	HealthzPort int
}

func RegisterService(input RegisterInput) error {
	for idx, ip := range envconf.EnvConf.IPs {
		reg := api.AgentServiceRegistration{
			ID:      fmt.Sprintf("%v-%v", input.ID, idx),
			Name:    input.Name,
			Tags:    input.Tags,
			Port:    input.Port,
			Address: ip,
		}

		chk := api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%v:%v/healthz", ip, input.Port),
			Timeout:                        "20s",
			Interval:                       "3s",
			DeregisterCriticalServiceAfter: "60s",
		}

		if envconf.EnvConf.ContainerID != envconf.NotRunInContainer {
			chk.DockerContainerID = envconf.EnvConf.ContainerID
		}

		reg.Check = &chk

		err := myClient.Agent().ServiceRegister(&reg)
		if err != nil {
			return xerrors.Errorf("fail to register service for %v: %v", ip, err)
		}
	}

	return nil
}

func DeregisterService(id uuid.UUID) error {
	return myClient.Agent().ServiceDeregister(fmt.Sprintf("%v", id))
}

func QueryServices(serviceName string) (map[string]*api.AgentService, error) {
	return myClient.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%v\"", serviceName))
}
