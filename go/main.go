package main

import (
	"log/slog"
	"os"
	{{ if .gin -}}
	"context"
	"time"
	{{- end }}

	"github.com/spf13/cobra"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
	{{ if .gin -}}
	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/router"
	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/server"
	{{- end }}
)

// Main function that will be executed from the root command.
func run(rootCmd *cobra.Command) {
	conf, err := cmd.NewConf(rootCmd)
	if err != nil {
		slog.Error("unable to load configuration", "error", err)
		os.Exit(1)
	}

	lg := cmd.NewLogger(conf)
	slog.SetDefault(lg)
	lg.Info("starting", "version", cmd.Version, "build", cmd.Build, "date", cmd.BuildDate)
	{{- if .gin }}

	tel, err := server.NewTelemetry(conf, lg)
	if err != nil {
		lg.Error("unable to setup telemetry", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tel.Shutdown(ctx); err != nil {
			lg.Error("telemetry shutdown error", "error", err)
		}
	}()

	// Swap in combined logger (stderr + OTLP when instrumented)
	lg = tel.Logger(lg)
	slog.SetDefault(lg)

	cc, err := server.NewCors(conf, lg)
	if err != nil {
		lg.Error("unable to setup cors", "error", err)
		os.Exit(1)
	}
	e := server.NewGinEngine(conf, lg, cc, tel)
	router.New(conf, lg, e).Listen()
	{{- end }}
}

func main() {
	root := &cobra.Command{
		Use:     "{{ .name }}",
		Short:   "{{ .description }}",
		Version: cmd.Version,
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd)
		},
	}

	cmd.Setup(root)

	// Run the command
	if err := root.Execute(); err != nil {
		slog.Error("unable to start", "error", err)
		os.Exit(1)
	}
}
