package dtm

import (
	"fmt"
	"golang.org/x/xerrors"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const"
	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const"
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

func GetGrpcUrl(serviceName string, grpcName string, service ...string) (string, error) {

	grpcUrl, err := SetPackageAndService(serviceName, service...)
	if err != nil {
		return "", err
	}
	grpcUrl = grpcUrl + "/" + grpcName
	return grpcUrl, nil
}

func SetPackageAndService(serviceName string, service ...string) (string, error) {
	switch serviceName {
	case appusermgrconst.ServiceName:
		serviceName, err := GetService(serviceName)
		if err != nil {
			return "", err
		}
		if len(service) != 0 {
			return serviceName + "/app.user.manager.v1." + service[0], nil
		}
		return serviceName + "/app.user.manager.v1.AppUserManager", nil
	case inspireconst.ServiceName:
		serviceName, err := GetService(serviceName)
		if err != nil {
			return "", err
		}
		if len(service) != 0 {
			return serviceName + "/cloud.hashing.inspire.v1." + service[0], nil
		}
		return serviceName + "/cloud.hashing.inspire.v1.CloudHashingInspire", nil
	}
	return serviceName, nil
}
