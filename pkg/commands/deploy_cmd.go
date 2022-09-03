package commands

import (
	"fmt"

	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredDeployFlags = []string{"clusterName", "serviceName", "image"}

// NewDeployCmd creates cli command for deploying new version of ecs service
func NewDeployCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "deploy",
		Aliases: []string{"d"},
		Usage:   "Deploys new version of app, takes newest task definition and updates docker image by creating new revision",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "clusterName",
				Usage: "Name of ECS cluster to which new version should be deployed",
			},
			cli.StringFlag{
				Name:  "serviceName",
				Usage: "existing service in ECS cluster which should be updated",
			},
			cli.StringSliceFlag{
				Name:  "image",
				Usage: "Image with tag which will be used to create new Task Definition",
			},
			cli.IntSliceFlag{
				Name:  "imageIndex",
				Usage: "Index of image in container definitions that should be updated",
				Value: &cli.IntSlice{},
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
			images := c.StringSlice("image")
			imageIndexes := c.IntSlice("imageIndex")

			if len(imageIndexes) == 0 {
				imageIndexes = append(imageIndexes, 0)
			}

			if len(imageIndexes) != len(images) {
				return fmt.Errorf("imageIndexes and images must be repeated same amount of times")
			}

			codedeployGroup := c.String("codedeployGroup")
			codedeployApp := c.String("codedeployApp")

			if len(codedeployGroup) == 0 {
				codedeployGroup = serviceName
			}

			if len(codedeployApp) == 0 {
				codedeployApp = serviceName
			}

			output, err := deployment.Deploy(&clusterName, &serviceName, images, imageIndexes, &codedeployApp, &codedeployGroup)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
