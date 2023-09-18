package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type LogConf struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Caller bool   `mapstructure:"caller"`
}
{{ if .gin }}
type ServerConf struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Mode       string `mapstructure:"mode"`
	Instrument bool   `mapstructure:"instrument"`

	Cors CorsConf `mapstructure:"cors"`
}

type CorsConf struct {
	Enabled bool     `mapstructure:"enabled"`
	Methods []string `mapstructure:"methods"`
	Headers []string `mapstructure:"headers"`
	Expose  []string `mapstructure:"expose"`
	Origins []string `mapstructure:"origins"`
	All     bool     `mapstructure:"all"`
}
{{ end }}
type Conf struct {
	Log    LogConf    `mapstructure:"log"`
	{{- if .gin }}
	Server ServerConf `mapstructure:"server"`
	{{- end }}
}

// NewLogger will return a new logger
func NewLogger(c *Conf) zerolog.Logger {
	// Level parsing
	warns := []string{}
	lvl, err := zerolog.ParseLevel(c.Log.Level)
	if err != nil {
		warns = append(warns, fmt.Sprintf("unrecognized log level '%s', fallback to 'info'", c.Log.Level))
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(lvl)
	}

	// Type parsing
	switch c.Log.Format {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "json":
		break
	default:
		warns = append(warns, fmt.Sprintf("unrecognized log format '%s', fallback to 'json'", c.Log.Format))
	}

	// Caller
	if c.Log.Caller {
		log.Logger = log.With().Caller().Logger()
	}

	// Log messages with the newly created logger
	for _, w := range warns {
		log.Warn().Msg(w)
	}

	return log.Logger
}

// NewConf will parse and return the configuration
func NewConf() (*Conf, error) {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("project-name")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Configuration file
	if viper.GetString("conf") != "" {
		viper.SetConfigFile(viper.GetString("conf"))
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config/")
	}

	viper.ReadInConfig() // nolint: errcheck
	conf := &Conf{}
	if err := viper.Unmarshal(conf); err != nil {
		return conf, fmt.Errorf("unable to unmarshal conf: %w", err)
	}

	return conf, nil
}
