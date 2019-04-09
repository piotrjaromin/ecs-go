package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredListDeploymentsFlags = []string{"codedeployApp", "codedeployGroup"}

func NewListDeploymentsCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "list-deployments",
		Aliases: []string{"ld"},
		Usage:   "Deploys new version of app",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "codedeployApp",
				Usage: "codedeploy application for which deployments should be listed",
			},
			cli.StringFlag{
				Name:  "codedeployGroup",
				Usage: "codedeployGroup group for which deployments should be listed",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredListDeploymentsFlags); err != nil {
				return err
			}

			codedeployApp := c.String("codedeployApp")
			codedeployGroup := c.String("codedeployGroup")
			output, err := deployment.ListDeployments(&codedeployApp, &codedeployGroup)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
