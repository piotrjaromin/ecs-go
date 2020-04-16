package commands

import (
	"github.com/piotrjaromin/ecs-go/pkg/services"
	"github.com/urfave/cli"
)

var requiredTagImageFlags = []string{"repositoryName", "currentTag", "newTag"}

// NewTagImageCmd creates cli command for adding new tag to image identified by currentTag
func NewTagImageCmd(deployment services.Deployment) cli.Command {
	return cli.Command{
		Name:    "tag-image",
		Aliases: []string{"ti"},
		Usage:   "Tags ECR image",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "repositoryName",
				Usage: "Name of ECR repositoryd",
			},
			cli.StringFlag{
				Name:  "currentTag",
				Usage: "existing tag used to search the image",
			},
			cli.StringFlag{
				Name:  "newTag",
				Usage: "new tag to be set on the image",
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredTagImageFlags); err != nil {
				return err
			}

			repositoryName := c.String("repositoryName")
			currentTag := c.String("currentTag")
			newTag := c.String("newTag")

			err := deployment.TagImage(&repositoryName, &currentTag, &newTag)
			if err != nil {
				return err
			}

			return printOutput(true)
		},
	}
}
