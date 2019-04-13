package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredWaitForStateFlags = []string{"deploymentId", "state"}

func NewWaitForStateCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "wait-for-state",
		Aliases: []string{"cd"},
		Usage:   "waits until given deployment reaches given state",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "deploymentId",
				Usage: "deployment id for which we are waiting",
			},
			cli.StringFlag{
				Name:  "state",
				Usage: "state which given deployment should reach",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredWaitForStateFlags); err != nil {
				return err
			}

			// deploymentID := c.String("deploymentId")
			// state := c.String("state")

			return nil
			// output, err := deployment.RollbackLatestDeployment(&codedeployApp, &codedeployGroup)
			// if err != nil {
			// 	return err
			// }

			// return printOutput(output)
		},
	}
}
