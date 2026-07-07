package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"{{.gitserver}}/{{.owner}}/{{.name}}/cmd"
	"{{.gitserver}}/{{.owner}}/{{.name}}/router"

	_ "{{.gitserver}}/{{.owner}}/{{.name}}/migrations"
)

func main() {
	app := pocketbase.New()

	// Register custom commands and flags
	app.RootCmd.AddCommand(cmd.VersionCmd)
	app.RootCmd.PersistentFlags().StringP("conf", "c", "", "path to configuration file (also settable via {{.name | uc}}_CONF env var)")

	// Load configuration
	c, err := cmd.NewConf()
	if err != nil {
		log.Fatal(err)
	}

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: c.Dev,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		return router.Setup(se)
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
