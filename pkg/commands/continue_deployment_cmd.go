package commands

import (
	"fmt"
	"github.com/urfave/cli"
)

func NewContinueDeploymentsCmd() cli.Command {
	return cli.Command{
		Name:    "continue-deployment",
		Aliases: []string{"cd"},
		Usage:   "Allows active deployment to continue deployment",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "deploymentId",
				Usage: "Id of deployment which should be continued",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("creating deployment: ", c.Args().First())
			return nil
		},
	}
}
