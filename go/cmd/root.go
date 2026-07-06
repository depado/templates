package cmd

import (
	"github.com/spf13/cobra"
)

// Setup takes the root command, adds persistent flags and subcommands.
// Flag binding to viper is handled in NewConf to keep a local viper instance.
func Setup(root *cobra.Command) {
	addConfigurationFlag(root)
	addLoggerFlags(root)
	{{ if .gin -}}
	addServerFlags(root)
	{{- end }}

	root.AddCommand(versionCmd)
}
