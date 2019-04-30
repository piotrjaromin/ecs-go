package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredScaleFlags = []string{"clusterName", "serviceName", "count"}

// NewScaleCmd creates cli command for sclaing number of tasks of ecs service
func NewScaleCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "scale",
		Aliases: []string{"s"},
		Usage:   "Changes number of tasks",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "clusterName",
				Usage: "Name of ECS cluster in which service is deployed",
			},
			cli.StringFlag{
				Name:  "serviceName",
				Usage: "existing service in ECS cluster which should be updated",
			},
			cli.UintFlag{
				Name:  "count",
				Usage: "New value for service instance count",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredScaleFlags); err != nil {
				return err
			}

			clusterName := c.String("clusterName")
			serviceName := c.String("serviceName")
			count := c.Uint("count")

			output, err := deployment.Scale(&clusterName, &serviceName, count)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
