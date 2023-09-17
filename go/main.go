package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// Main function that will be executed from the root command.
func run() {
	conf, err := cmd.NewConf()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load configuration")
	}

	lg := cmd.NewLogger(conf)
	lg.Info().Msg("hello world")
}

func main() {
	// Initialize Cobra and Viper
	rootCmd := &cobra.Command{
		Use:     "{{ .name }}",
		Short:   "{{ .description }}",
		Version: cmd.Version,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	cmd.AddAllFlags(rootCmd)
	rootCmd.AddCommand(cmd.VersionCmd)

	// Run the command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("unable to start")
	}
}
