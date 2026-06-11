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
	c, err := cmd.NewConf()
	if err != nil {
		log.Fatal(err)
	}

	app := pocketbase.New()
	app.RootCmd.AddCommand(cmd.VersionCmd)

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
