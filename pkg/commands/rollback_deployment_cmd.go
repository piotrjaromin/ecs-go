package commands

import (
	"fmt"
	"github.com/urfave/cli"
)

func NewRollbackDeploymentCmd() cli.Command {
	return cli.Command{
		Name:    "rollback-deployment",
		Aliases: []string{"rd"},
		Usage:   "Rollbacks active deployment",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "deploymentId",
				Usage: "Id of deployment which should be rolledback",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("creating deployment: ", c.Args().First())
			return nil
		},
	}
}
