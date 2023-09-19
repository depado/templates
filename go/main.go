package main

import (
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).Msg("unable to load configuration")
	}

	lg := cmd.NewLogger(conf)
	lg.Info().Str("version", cmd.Version).Str("build", cmd.Build).Str("date", cmd.Time).Send()
	{{- if .gin }}

	cc, err := server.NewCors(conf, lg)
	if err != nil {
		lg.Fatal().Err(err).Msg("unable to setup cors")
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
		log.Fatal().Err(err).Msg("unable to start")
	}
}
