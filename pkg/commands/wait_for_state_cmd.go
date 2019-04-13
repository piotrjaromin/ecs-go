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
			cli.IntFlag{
				Name:  "waitTime",
				Usage: "max time in seconds after which this command will end, defaults to 30 minutes",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredWaitForStateFlags); err != nil {
				return err
			}

			deploymentID := c.String("deploymentId")
			state := c.String("state")
			waitTime := c.Int("waitTime")
			if waitTime == 0 {
				waitTime = 30 * 60 // by default wait 30 minutes
			}

			output, err := deployment.WaitForState(&deploymentID, &state, waitTime)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
