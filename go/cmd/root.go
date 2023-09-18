package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Setup takes the root command, binds the flags to it and adds the other
// commands.
func Setup(root *cobra.Command) {
	addConfigurationFlag(root)
	addLoggerFlags(root)
	{{ if .gin -}}
	addServerFlags(root)
	{{- end }}

	if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
		log.Fatal().Err(err).Msg("unable to bind flags")
	}

	root.AddCommand(versionCmd)

}
