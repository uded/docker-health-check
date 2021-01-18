// +build http

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app.Action = func(c *cli.Context) error {
		if resp, err := httpResp(c); err != nil {
			return err
		} else if resp.StatusCode != statusCode {
			return cli.Exit(fmt.Sprintf("resp code %d didn't match %d", resp.StatusCode, statusCode), 1)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
