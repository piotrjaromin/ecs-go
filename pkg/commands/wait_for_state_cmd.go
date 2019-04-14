package commands

import (
	"fmt"
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
	"strings"
)

var requiredWaitForStateFlags = []string{"deploymentId"}

var validStates = []string{"InProgress", "Ready", "Created", "Queued", "Stopped", "Failed", "Succeeded"}

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
				Usage: "state which given deployment should reach, by deafult it is 'Ready'",
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

			if len(state) == 0 {
				state = "Ready"
			}

			if !hasString(state, validStates) {
				return fmt.Errorf("Invalid state provided %s, valid values are: %s", state, strings.Join(validStates, ", "))
			}

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

func hasString(searched string, array []string) bool {
	for _, current := range array {
		if searched == current {
			return true
		}
	}

	return false
}
