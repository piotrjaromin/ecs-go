package commands

import (
	"fmt"
	"github.com/urfave/cli"
)

func NewListDeploymentsCmd() cli.Command {
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
			fmt.Println("creating deployment: ", c.Args().First())
			return nil
		},
	}
}
