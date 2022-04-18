package dtm

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
)

func GetService(serviceName ...string) (string, error) {
	service := "dtm-cluster.npool.top"
	if len(serviceName) != 0 {
		service = serviceName[0]
		svc, err := config.PeekService(service, grpc2.GRPCTAG)
		if err != nil {
			return "", xerrors.Errorf("Fail to query %v service: %v", service, err)
		}
		return fmt.Sprintf("%v:%v", svc.Address, svc.Port), nil
	}
	svc, err := config.PeekService(service)
	if err != nil {
		return "", xerrors.Errorf("Fail to query %v service: %v", service, err)
	}
	return fmt.Sprintf("%v:%v", svc.Address, svc.Port), nil
}
