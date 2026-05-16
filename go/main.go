package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
	{{ if .gin -}}
	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/router"
	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/server"
	{{- end }}
)

// Main function that will be executed from the root command.
func run() {
	conf, err := cmd.NewConf()
	if err != nil {
		slog.Error("unable to load configuration", "error", err)
		os.Exit(1)
	}

	lg := cmd.NewLogger(conf)
	slog.SetDefault(lg)
	lg.Info("starting", "version", cmd.Version, "build", cmd.Build, "date", cmd.BuildDate)
	{{- if .gin }}

	cc, err := server.NewCors(conf, lg)
	if err != nil {
		lg.Error("unable to setup cors", "error", err)
		os.Exit(1)
	}
	e := server.NewGinEngine(conf, lg, cc)
	router.New(conf, lg, e).Listen()
	{{- end }}
}

func main() {
	root := &cobra.Command{
		Use:     "{{ .name }}",
		Short:   "{{ .description }}",
		Version: cmd.Version,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	cmd.Setup(root)

	// Run the command
	if err := root.Execute(); err != nil {
		slog.Error("unable to start", "error", err)
		os.Exit(1)
	}
}
