package main

import (
	"log"
	"os"

	"github.com/piotrjaromin/ecs-go/pkg/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		commands.NewDeployCmd(),
		commands.NewContinueDeploymentsCmd(),
		commands.NewListDeploymentsCmd(),
		commands.NewRollbackDeploymentCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
