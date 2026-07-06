package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Conf holds the application configuration.
type Conf struct {
	// Dev enables development mode (automigrate, debug logging, etc.)
	Dev bool `mapstructure:"dev"`

	// Add your configuration fields here, for example:
	// API APIConf `mapstructure:"api"`
}

// NewConf will parse and return the configuration.
// The config file path can be set via the APPNAME_CONF environment variable.
func NewConf() (*Conf, error) {
	v := viper.New()

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("{{ .name }}")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configuration file
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
