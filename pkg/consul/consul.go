package consul

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	*api.Client
	envConf *envconf.EnvConf
}

func NewConsulClient() (*Client, error) {
	envConf, err := envconf.NewEnvConf()
	if err != nil {
		return nil, xerrors.Errorf("fail to create environment configuration: %v", err)
	}

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%v:%v", envConf.ConsulHost, envConf.ConsulPort)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, xerrors.Errorf("fail to create consul client: %v", err)
	}

	return &Client{
		Client:  client,
		envConf: envConf,
	}, nil
}

// IP is parsed from package envconf
type RegisterInput struct {
	ID          uuid.UUID
	Name        string
	Tags        []string
	Port        int
	HealthzPort int
}

func (c *Client) RegisterService(input RegisterInput) error {
	for idx, ip := range c.envConf.IPs {
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

		if c.envConf.ContainerID != envconf.NotRunInContainer {
			chk.DockerContainerID = c.envConf.ContainerID
		}

		reg.Check = &chk

		fmt.Printf("register service for %v\n", ip)

		err := c.Agent().ServiceRegister(&reg)
		if err != nil {
			return xerrors.Errorf("fail to register service for %v: %v", ip, err)
		}
	}

	return nil
}

func (c *Client) DeregisterService(id uuid.UUID) error {
	return c.Agent().ServiceDeregister(fmt.Sprintf("%v", id))
}

func (c *Client) QueryServices(serviceName string) (map[string]*api.AgentService, error) {
	return c.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%v\"", serviceName))
}
