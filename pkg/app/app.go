package app

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/version"
	banner "github.com/common-nighthawk/go-figure"
	cli "github.com/urfave/cli/v2"
)

var versionCmd = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "print version",
	Action: func(ctx *cli.Context) error {
		ver, err := version.GetVersion()
		if err != nil {
			return err
		}
		logger.Sugar().Infow(ver)
		return nil
	},
}

func NewApp(
	serviceName, description, usageText, argsUsage string,
	flags []cli.Flag,
	authors []*cli.Author,
	commands []*cli.Command) *cli.App {
	banner.NewColorFigure(serviceName, "", "green", true).Print()
	commands = append(commands, versionCmd)

	app := &cli.App{
		Name:        serviceName,
		Version:     "0.1.0",
		Description: description,
		ArgsUsage:   argsUsage,
		Usage:       usageText,
		Flags:       flags,
		Commands:    commands,
	}

	return app
}
