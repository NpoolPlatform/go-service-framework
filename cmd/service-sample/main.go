package main

import (
	"os"

	"github.com/NpoolPlatform/go-service-framework/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	_ "github.com/urfave/cli/v2"
)

const serviceName = "Service Sample"

func main() {
	err := app.NewApp(serviceName, "", "", "", nil, nil, nil).Run(os.Args)
	if err != nil {
		logger.Sugar().Infof("fail to run %v: %v", serviceName, err)
	}
}
