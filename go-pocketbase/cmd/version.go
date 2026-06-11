package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build number and versions injected at compile time via ldflags.
var (
	Version   = "unknown"
	Build     = "unknown"
	BuildDate = "unknown"
)

var versionHelp = `
This command will output the build number, version number and build date of {{ .name }}.
The build number corresponds to the sha1 commit the binary was built against,
while the version number corresponds to the latest tag the binary was built on.
Finally the build date corresponds to the date the binary was built.

If both values are "unknown" make sure to build {{ .name }} with "make".
`

// VersionCmd displays the build number, version and build date.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build, version and build date",
	Long:  versionHelp,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build: %s\nVersion: %s\nBuild Date: %s\n", Build, Version, BuildDate)
	},
}
