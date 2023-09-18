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
	{{ if not .gin -}}
	lg.Info().Msg("hello world")
	{{- end }}
	{{ if .gin -}}
	cc, err := server.NewCors(conf, lg)
	if err != nil {
		lg.Fatal().Err(err).Msg("unable to setup cors")
	}
	e := server.NewGinEngine(conf, lg, cc)
	r := router.New(conf, lg, e)
	if err = r.Listen(); err != nil {
		lg.Fatal().Err(err).Msg("unable to run")
	}
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
