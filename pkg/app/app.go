package app

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/version"
	banner "github.com/common-nighthawk/go-figure"
	cli "github.com/urfave/cli/v2"
)

func NewApp(
	serviceName, description, usageText, argsUsage string,
	flags []cli.Flag,
	authors []*cli.Author,
	commands []*cli.Command) (*cli.App, error) {
	banner.NewColorFigure(serviceName, "", "green", true).Print()

	ver, err := version.GetVersion()
	if err != nil {
		return nil, err
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

	return app, nil
}
