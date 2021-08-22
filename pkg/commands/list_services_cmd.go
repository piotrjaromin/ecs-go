package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

// NewListServicesCmd creates cli command for listing all ECS services from current region
func NewListServicesCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "list-services",
		Aliases: []string{"ls"},
		Usage:   "Gets list of services",
		Action: func(c *cli.Context) error {
			output, err := deployment.ListServices()
			if err != nil {
				return err
			}

			return printOutput(output)
		},
	}
}
