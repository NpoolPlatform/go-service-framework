package consul

import (
	"fmt"
	"strings"

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

func RegisterService(checkHealth bool, input RegisterInput) error {
	addr := serviceName2PODService(input.Name)
	reg := api.AgentServiceRegistration{
		ID:      input.ID.String(),
		Name:    input.Name,
		Tags:    input.Tags,
		Port:    input.Port,
		Address: addr,
	}

	if checkHealth {
		chk := api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%v:%v/healthz", addr, input.Port),
			Timeout:                        "20s",
			Interval:                       "3s",
			DeregisterCriticalServiceAfter: "60s",
		}

		if envconf.EnvConf.ContainerID != envconf.NotRunInContainer {
			chk.DockerContainerID = envconf.EnvConf.ContainerID
		}

		reg.Check = &chk
	}

	err := myClient.Agent().ServiceRegister(&reg)
	if err != nil {
		return xerrors.Errorf("fail to register service for %v: %v", addr, err)
	}

	return nil
}

func serviceName2PODService(name string) string {
	topDomain := "npool.top"
	k8sNS := "kube-system"
	k8sDNS := "svc.cluster.local"
	return fmt.Sprintf("%s%s.%s", strings.TrimSuffix(name, topDomain), k8sNS, k8sDNS)
}

func DeregisterService(id uuid.UUID) error {
	return myClient.Agent().ServiceDeregister(fmt.Sprintf("%v", id))
}

func QueryServices(serviceName string, tags ...string) (map[string]*api.AgentService, error) {
	configs, err := myClient.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%v\"", serviceName))
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return configs, nil
	}

	tagm := func() map[string]struct{} {
		tm := make(map[string]struct{})
		for _, tag := range tags {
			tm[tag] = struct{}{}
		}
		return tm
	}()

	cfgs := make(map[string]*api.AgentService)
	for k, as := range configs {
		for _, t := range as.Tags {
			if _, ok := tagm[t]; ok {
				cfgs[k] = as
			}
		}
	}

	return cfgs, nil
}
