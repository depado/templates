package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/router"

	_ "{{ .gitserver }}/{{ .owner }}/{{ .name }}/migrations"
)

func main() {
	app := pocketbase.New()

	// Register custom commands
	app.RootCmd.AddCommand(cmd.VersionCmd)

	// Load configuration
	c, err := cmd.NewConf()
	if err != nil {
		log.Fatal(err)
	}

	// Register migrate command with automigrate in dev mode
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: c.Dev,
	})

	// Setup routes and hooks before serving
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		return router.Setup(se)
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
