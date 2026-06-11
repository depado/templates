package cmd

import (
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
func NewConf() (*Conf, error) {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("{{ .name }}")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configuration file
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
