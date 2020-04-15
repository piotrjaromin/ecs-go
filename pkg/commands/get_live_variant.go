package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredGetLiveVariantFlags = []string{"clusterName", "serviceName"}

// NewGetLiveVariantCmd creates cli command for resolving current live variant of ecs service (blue|green)
func NewGetLiveVariantCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "get-live-variant",
		Aliases: []string{"glv"},
		Usage:   "Gets live variant of service",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "clusterName",
				Usage: "Name of ECS cluster in which service is deployed",
			},
			cli.StringFlag{
				Name:  "serviceName",
				Usage: "existing service in ECS cluster",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredGetLiveVariantFlags); err != nil {
				return err
			}

			clusterName := c.String("clusterName")
			serviceName := c.String("serviceName")

			output, err := deployment.GetLiveVariant(&clusterName, &serviceName)
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
