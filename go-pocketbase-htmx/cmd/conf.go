package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Conf struct {
	Dev bool `mapstructure:"dev"`
}

func NewConf() (*Conf, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("{{.name}}")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if viper.GetString("conf") != "" {
		viper.SetConfigFile(viper.GetString("conf"))
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config/")
	}

	viper.ReadInConfig() //nolint:errcheck

	conf := &Conf{}
	if err := viper.Unmarshal(conf); err != nil {
		return conf, fmt.Errorf("unable to unmarshal conf: %w", err)
	}

	return conf, nil
}
