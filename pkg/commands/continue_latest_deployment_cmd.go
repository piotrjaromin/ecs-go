package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredContinueLatestDeployFlags = []string{"codedeployApp", "codedeployGroup"}

func NewContinueLatestDeploymentCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "continue-latest-deployment",
		Aliases: []string{"cd"},
		Usage:   "Trigger continue on latest deployment for given app and group",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "codedeployApp",
				Usage: "codedeploy application which is used to continue deployment",
			},
			cli.StringFlag{
				Name:  "codedeployGroup",
				Usage: "codedeployGroup group which is used to continue deployment",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredContinueLatestDeployFlags); err != nil {
				return err
			}

			codedeployGroup := c.String("codedeployGroup")
			codedeployApp := c.String("codedeployApp")

			output, err := deployment.ContinueLatestDeployment(&codedeployApp, &codedeployGroup)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
