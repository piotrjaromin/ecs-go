package commands

import (
	"fmt"
	"github.com/urfave/cli"
)

var requiredDeployFlags = []string{"clusterName", "serviceName", "image"}

// NewDeployCmd creates cli command for deploying new version of ecs service
func NewDeployCmd() cli.Command {
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

			fmt.Println("creating deployment: ", c.Args().First())
			return nil
		},
	}
}
