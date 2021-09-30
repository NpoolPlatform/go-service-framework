package config

import (
	"flag"
	"fmt"

	"golang.org/x/xerrors"

	"github.com/philchia/agollo/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	agollo.Client
}

func Init(configPath, appName string) (*Config, error) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, xerrors.Errorf("fail to bind flags: %v", err)
	}

	viper.SetConfigName(fmt.Sprintf("%s.viper", appName))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(fmt.Sprintf("/etc/%v", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%v", appName))
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, xerrors.Errorf("fail to init config: %v", err)
	}

	apolloCfg := viper.GetStringMap("apollo")

	names := make([]string, 0)
	for _, namespacename := range apolloCfg["namespacenames"].([]interface{}) {
		names = append(names, namespacename.(string))
	}

	agolloCli := agollo.NewClient(&agollo.Conf{
		AppID:          apolloCfg["appid"].(string),
		Cluster:        apolloCfg["cluster"].(string),
		NameSpaceNames: names,
		MetaAddr:       apolloCfg["metaaddr"].(string),
	})

	err = agolloCli.Start()
	if err != nil {
		return nil, xerrors.Errorf("fail to start apollo client: %v", err)
	}

	return &Config{
		Client: agolloCli,
	}, nil
}
