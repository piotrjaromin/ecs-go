package commands

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
)

func validateRequiredFlags(c *cli.Context, flags []string) error {
	errMsg := ""

	for _, f := range flags {
		if len(c.String(f)) == 0 {
			errMsg = fmt.Sprintf("%s\nMissing required flag --%s", errMsg, f)
		}
	}

	if len(errMsg) > 0 {
		return fmt.Errorf(errMsg)
	}

	return nil
}

func printOutput(output interface{}) error {
	data, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)
	return nil
}
