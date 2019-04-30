package main

import (
	"log"
	"os"

	"github.com/piotrjaromin/ecs-go/pkg/commands"
	"github.com/urfave/cli"

	"github.com/piotrjaromin/ecs-go/pkg/services"
)

func main() {
	app := cli.NewApp()

	deploySvc, err := services.NewDeployment()

	if err != nil {
		panic(err)
	}

	app.Commands = []cli.Command{
		commands.NewDeployCmd(deploySvc),
		commands.NewScaleCmd(deploySvc),
		commands.NewContinueDeploymentsCmd(deploySvc),
		commands.NewForceContinueDeploymentsCmd(deploySvc),
		commands.NewListDeploymentsCmd(deploySvc),
		commands.NewRollbackDeploymentCmd(deploySvc),

		commands.NewWaitForStateCmd(deploySvc),
		commands.NewContinueLatestDeploymentCmd(deploySvc),
		commands.NewForceContinueLatestDeploymentCmd(deploySvc),
		commands.NewRollbackLatestDeploymentCmd(deploySvc),
		commands.NewWaitForLatestCmd(deploySvc),

		// TODO for Future
		// commands.CreateTaskDefinition(deploySvc),
		// commands.CreateECSService(deploySvc),
	}

	app.Name = "ecs integration for codedeploy"
	app.Description = "Cli application that trigger AWS ECS deployments through AWS Codedeploy service"
	app.Version = "1.0.0"
	app.Usage = ""

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
