package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredDeployFlags = []string{"clusterName", "serviceName", "image"}

// NewDeployCmd creates cli command for deploying new version of ecs service
func NewDeployCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "deploy",
		Aliases: []string{"d"},
		Usage:   "Deploys new version of app",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "clusterName",
				Usage: "Name of ECS cluster to which new version should be deployed",
			},
			cli.StringFlag{
				Name:  "serviceName",
				Usage: "existing service in ECS cluster which should be updated",
			},
			cli.StringFlag{
				Name:  "image",
				Usage: "Image with tag which will be used to create new Task Definition",
			},
			cli.StringFlag{
				Name:  "codedeployApp",
				Usage: "codedeploy application which is used to trigger deployment",
			},
			cli.StringFlag{
				Name:  "codedeployGroup",
				Usage: "codedeployGroup group which is used to trigger deployment",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredDeployFlags); err != nil {
				return err
			}

			clusterName := c.String("clusterName")
			serviceName := c.String("serviceName")
			image := c.String("image")

			codedeployGroup := c.String("codedeployGroup")
			codedeployApp := c.String("codedeployApp")

			if len(codedeployGroup) == 0 {
				codedeployGroup = serviceName
			}

			if len(codedeployApp) == 0 {
				codedeployApp = serviceName
			}

			output, err := deployment.Deploy(&clusterName, &serviceName, &image, &codedeployApp, &codedeployGroup)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
