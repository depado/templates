package cmd

import (
	"github.com/spf13/cobra"
)

// addLoggerFlags adds support to configure the level of the logger.
func addLoggerFlags(c *cobra.Command) {
	c.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	c.PersistentFlags().String("log.format", "console", `one of "console" or "json"`)
	c.PersistentFlags().Bool("log.caller", false, "display the file and line where the call was made")
}
{{ if .gin }}
// addServerFlags adds support to configure the server
func addServerFlags(c *cobra.Command) {
	// Server related flags
	c.PersistentFlags().String("server.host", "127.0.0.1", "host on which the server should listen")
	c.PersistentFlags().Int("server.port", 8080, "port on which the server should listen")
	c.PersistentFlags().String("server.mode", "release", "server mode can be either 'debug', 'test' or 'release'")
	c.PersistentFlags().Bool("server.instrument", true, "enable prometheus instrumentation")
	c.PersistentFlags().Bool("server.unified-logger", true, "use zerolog to log requests instead of gin's default logger")

	// CORS related flags
	c.PersistentFlags().Bool("server.cors.enabled", false, "enable CORS")
	c.PersistentFlags().StringSlice("server.cors.methods", []string{"GET", "PUT", "POST", "DELETE", "OPTION", "PATCH"}, "array of allowed method when cors is enabled")
	c.PersistentFlags().StringSlice("server.cors.headers", []string{"Origin", "Authorization", "Content-Type"}, "array of allowed headers")
	c.PersistentFlags().StringSlice("server.cors.expose", []string{}, "array of exposed headers")
	c.PersistentFlags().StringSlice("server.cors.origins", []string{}, "array of allowed origins (overwritten if all is active)")
	c.PersistentFlags().Bool("server.cors.all", false, "allow all origins (overrides origins if set)")
}
{{ end }}
// AddConfigurationFlag adds support to provide a configuration file on the
// command line.
func addConfigurationFlag(c *cobra.Command) {
	c.PersistentFlags().StringP("conf", "c", "", "configuration file to use")
}
