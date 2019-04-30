package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredForceContinueDeployFlags = []string{"deploymentId"}

func NewForceContinueDeploymentsCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "force-continue-deployment",
		Aliases: []string{"cd"},
		Usage:   "Forces active deployment to continue deployment (kills replacement task before its time elapsed",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "deploymentId",
				Usage: "Id of deployment which should be continued",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredContinueDeployFlags); err != nil {
				return err
			}

			deploymentID := c.String("deploymentId")
			output, err := deployment.ForceContinueDeployment(&deploymentID)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
