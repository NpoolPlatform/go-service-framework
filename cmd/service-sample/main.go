package main

import (
	"os"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	banner "github.com/common-nighthawk/go-figure"
	cli "github.com/urfave/cli/v2"
)

const serviceName = "Service Sample"

func main() {
	app := &cli.App{
		Name:  serviceName,
		Usage: "A sample service for all npool service framework",
		Action: func(c *cli.Context) error {
			banner.NewColorFigure(serviceName, "", "green", true).Print()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Sugar().Infof("fail to run %v: %v", serviceName, err)
	}
}
