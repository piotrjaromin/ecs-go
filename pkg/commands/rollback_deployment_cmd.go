package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredRollbackFlags = []string{"deploymentId"}

func NewRollbackDeploymentCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "rollback-deployment",
		Aliases: []string{"rd"},
		Usage:   "Rollbacks active deployment",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "deploymentId",
				Usage: "Id of deployment which should be rolledback",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredContinueDeployFlags); err != nil {
				return err
			}

			deploymentID := c.String("deploymentId")
			output, err := deployment.RollbackDeployment(&deploymentID)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
