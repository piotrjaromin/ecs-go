package commands

import (
	"fmt"
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
	"strings"
)

var requiredWaitForLatestFlags = []string{"codedeployApp", "codedeployGroup"}

func NewWaitForLatestCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "wait-for-latest",
		Aliases: []string{"cd"},
		Usage:   "waits until given deployment reaches given state",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "codedeployApp",
				Usage: "codedeploy application which is used to find latest deployment",
			},
			cli.StringFlag{
				Name:  "codedeployGroup",
				Usage: "codedeployGroup group which is used to find latest deployment",
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
			if err := validateRequiredFlags(c, requiredWaitForLatestFlags); err != nil {
				return err
			}

			codedeployApp := c.String("codedeployApp")
			codedeployGroup := c.String("codedeployGroup")
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

			output, err := deployment.WaitForLatest(&codedeployApp, &codedeployGroup, &state, waitTime)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
