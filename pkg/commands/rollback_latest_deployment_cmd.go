package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredRollbackLatestDeployFlags = []string{"codedeployApp", "codedeployGroup"}

func NewRollbackLatestDeploymentCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "rollback-latest-deployment",
		Aliases: []string{"cd"},
		Usage:   "Trigger rollback on latest deployment for given app and group",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "codedeployApp",
				Usage: "codedeploy application which is used to rollback deployment",
			},
			cli.StringFlag{
				Name:  "codedeployGroup",
				Usage: "codedeployGroup group which is used to rollback deployment",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredRollbackLatestDeployFlags); err != nil {
				return err
			}

			codedeployGroup := c.String("codedeployGroup")
			codedeployApp := c.String("codedeployApp")

			output, err := deployment.RollbackLatestDeployment(&codedeployApp, &codedeployGroup)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
