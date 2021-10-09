package app

import (
	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/mysql"
	"github.com/NpoolPlatform/go-service-framework/pkg/version"

	banner "github.com/common-nighthawk/go-figure"
	cli "github.com/urfave/cli/v2"
)

type App struct {
	app    *cli.App
	Config *config.Config
	Mysql  *mysql.Client
	Consul *consul.Client
}

var myApp App

func Init(
	serviceName, description, usageText, argsUsage string,
	flags []cli.Flag,
	authors []*cli.Author,
	commands []*cli.Command) error {
	banner.NewColorFigure(serviceName, "", "green", true).Print()

	ver, err := version.GetVersion()
	if err != nil {
		return xerrors.Errorf("Fail to get version: %v", err)
	}

	app := &cli.App{
		Name:        serviceName,
		Version:     ver,
		Description: description,
		ArgsUsage:   argsUsage,
		Usage:       usageText,
		Flags:       flags,
		Commands:    commands,
	}

	myApp.app = app
	myApp.Consul, err = consul.NewConsulClient()
	if err != nil {
		return xerrors.Errorf("Fail to create consul client: %v", err)
	}

	myApp.Config, err = config.Init("./", serviceName, myApp.Consul)
	if err != nil {
		return xerrors.Errorf("Fail to create configuration: %v", err)
	}

	return nil
}

func Run(args []string) error {
	return myApp.app.Run(args)
}
