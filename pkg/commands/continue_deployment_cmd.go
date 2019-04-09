package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredContinueDeployFlags = []string{"deploymentId"}

func NewContinueDeploymentsCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "continue-deployment",
		Aliases: []string{"cd"},
		Usage:   "Allows active deployment to continue deployment",
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
			output, err := deployment.ContinueDeployment(&deploymentID)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
