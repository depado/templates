package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Conf struct {
	Dev bool `mapstructure:"dev"`
}

// NewConf will parse and return the configuration.
// The config file path can be set via the APPNAME_CONF environment variable.
func NewConf() (*Conf, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("{{.name}}")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if v.GetString("conf") != "" {
		v.SetConfigFile(v.GetString("conf"))
	} else {
		v.SetConfigName("conf")
		v.AddConfigPath(".")
		v.AddConfigPath("/config/")
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	conf := &Conf{}
	if err := v.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("unable to unmarshal conf: %w", err)
	}

	return conf, nil
}
