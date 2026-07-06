package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LogConf struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Source bool   `mapstructure:"source"`
	Color  string `mapstructure:"color"`
}
{{ if .gin }}
type ServerConf struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Mode          string `mapstructure:"mode"`
	{{ if .gin_otel }}Instrument    bool   `mapstructure:"instrument"`{{ end }}
	UnifiedLogger bool   `mapstructure:"unified-logger"`

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

// NewLogger will return a new logger based on configuration
func NewLogger(c *Conf) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(c.Log.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
		slog.Warn("unrecognized log level, fallback to info", "level", c.Log.Level)
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: c.Log.Source,
	}

	var handler slog.Handler
	switch strings.ToLower(c.Log.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	case "text", "console":
		var noColor bool
		switch strings.ToLower(c.Log.Color) {
		case "always":
			noColor = false
		case "never":
			noColor = true
		default: // "auto" or empty
			noColor = !isatty.IsTerminal(os.Stderr.Fd())
		}
		handler = tint.NewHandler(os.Stderr, &tint.Options{
			Level:      level,
			AddSource:  c.Log.Source,
			TimeFormat: time.DateTime,
			NoColor:    noColor,
		})
	default:
		handler = slog.NewJSONHandler(os.Stderr, opts)
		slog.Warn("unrecognized log format, fallback to json", "format", c.Log.Format)
	}

	return slog.New(handler)
}

// NewConf will parse and return the configuration. It binds cobra flags to
// viper to ensure the correct precedence order: CLI flags > env vars > config
// file > defaults.
func NewConf(cmd *cobra.Command) (*Conf, error) {
	v := viper.New()

	if err := v.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, fmt.Errorf("unable to bind flags: %w", err)
	}

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
